package cluster

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/goextension/log"
	"github.com/ipfs/ipfs-cluster/utils"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	p2pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
)

// PriorityTag is used to attach metadata to peers in the peerstore
// so they can be sorted.
var PriorityTag = "cluster"

// Timeouts for network operations triggered by the PeerManager.
var (
	DNSTimeout     = 5 * time.Second
	ConnectTimeout = 5 * time.Second
)

// PeerManager provides utilities for handling cluster peer addresses
// and storing them in a libp2p Host peerstore.
type PeerManager struct {
	ctx           context.Context
	host          host.Host
	peerstoreLock sync.Mutex
	peerstorePath string
}

// New creates a PeerManager with the given libp2p Host and peerstorePath.
// The path indicates the place to persist and read peer addresses from.
// If empty, these operations (LoadPeerstore, SavePeerstore) will no-op.
func New(ctx context.Context, h host.Host, peerstorePath string) *PeerManager {
	return &PeerManager{
		ctx:           ctx,
		host:          h,
		peerstorePath: peerstorePath,
	}
}

// ImportPeer adds a new peer address to the host's peerstore, optionally
// dialing to it. The address is expected to include the /p2p/<peerID>
// protocol part or to be a /dnsaddr/multiaddress
// Peers are added with the given ttl.
func (p *PeerManager) ImportPeer(addr multiaddr.Multiaddr, connect bool, ttl time.Duration) (peer.ID, error) {
	if p.host == nil {
		return "", nil
	}

	protos := addr.Protocols()
	if len(protos) > 0 && protos[0].Code == madns.DnsaddrProtocol.Code {
		// We need to pre-resolve this
		log.Debugf("resolving %s", addr)
		ctx, cancel := context.WithTimeout(p.ctx, DNSTimeout)
		defer cancel()

		resolvedAddrs, err := madns.Resolve(ctx, addr)
		if err != nil {
			return "", err
		}
		if len(resolvedAddrs) == 0 {
			return "", fmt.Errorf("%s: no resolved addresses", addr)
		}
		var pid peer.ID
		for _, add := range resolvedAddrs {
			pid, err = p.ImportPeer(add, connect, ttl)
			if err != nil {
				return "", err
			}
		}
		return pid, nil // returns the last peer ID
	}

	pinfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return "", err
	}

	// Do not add ourselves
	if pinfo.ID == p.host.ID() {
		return pinfo.ID, nil
	}

	log.Debugf("adding peer address %s", addr)
	p.host.Peerstore().AddAddrs(pinfo.ID, pinfo.Addrs, ttl)

	if connect {
		go func() {
			ctx, cancel := context.WithTimeout(p.ctx, ConnectTimeout)
			defer cancel()
			p.host.Connect(ctx, *pinfo)
		}()
	}
	return pinfo.ID, nil
}

// RmPeer clear all addresses for a given peer ID from the host's peerstore.
func (p *PeerManager) RmPeer(pid peer.ID) error {
	if p.host == nil {
		return nil
	}

	log.Debugf("forgetting peer %s", pid.Pretty())
	p.host.Peerstore().ClearAddrs(pid)
	return nil
}

// if the peer has dns addresses, return only those, otherwise
// return all.
func (p *PeerManager) filteredPeerAddresses(id peer.ID) []multiaddr.Multiaddr {
	all := p.host.Peerstore().Addrs(id)
	var pas []multiaddr.Multiaddr
	var pdnsas []multiaddr.Multiaddr

	for _, a := range all {
		if madns.Matches(a) {
			pdnsas = append(pdnsas, a)
		} else {
			pas = append(pas, a)
		}
	}

	if len(pdnsas) > 0 {
		return pdnsas
	}

	sort.Sort(utils.ByString(pas))
	return pas
}

// PeerInfos returns a slice of peerinfos for the given set of peers in order
// of priority. For peers for which we know DNS
// multiaddresses, we only include those. Otherwise, the AddrInfo includes all
// the multiaddresses known for that peer. Peers without addresses are not
// included.
func (p *PeerManager) PeerInfos(peers []peer.ID) []peer.AddrInfo {
	if p.host == nil {
		return nil
	}

	if peers == nil {
		return nil
	}

	var pinfos []peer.AddrInfo
	for _, pr := range peers {
		if pr == p.host.ID() {
			continue
		}
		pinfo := peer.AddrInfo{
			ID:    pr,
			Addrs: p.filteredPeerAddresses(pr),
		}
		if len(pinfo.Addrs) > 0 {
			pinfos = append(pinfos, pinfo)
		}
	}

	toSort := &peerSort{
		pinfos: pinfos,
		pstore: p.host.Peerstore(),
	}
	// Sort from highest to lowest priority
	sort.Sort(toSort)

	return toSort.pinfos
}

// ImportPeers calls ImportPeer for every address in the given slice, using the
// given connect parameter. Peers are tagged with priority as given
// by their position in the list.
func (p *PeerManager) ImportPeers(addrs []multiaddr.Multiaddr, connect bool, ttl time.Duration) error {
	for i, a := range addrs {
		pid, err := p.ImportPeer(a, connect, ttl)
		if err == nil {
			p.SetPriority(pid, i)
		}
	}
	return nil
}

// ImportPeersFromPeerstore reads the peerstore file and calls ImportPeers with
// the addresses obtained from it.
func (p *PeerManager) ImportPeersFromPeerstore(connect bool, ttl time.Duration) error {
	return p.ImportPeers(p.LoadPeerstore(), connect, ttl)
}

// LoadPeerstore parses the peerstore file and returns the list
// of addresses read from it.
func (p *PeerManager) LoadPeerstore() (addrs []multiaddr.Multiaddr) {
	if p.peerstorePath == "" {
		return
	}
	p.peerstoreLock.Lock()
	defer p.peerstoreLock.Unlock()

	f, err := os.Open(p.peerstorePath)
	if err != nil {
		return // nothing to load
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addrStr := scanner.Text()
		if len(addrStr) == 0 || addrStr[0] != '/' {
			// skip anything that is not going to be a multiaddress
			continue
		}
		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			log.Errorf(
				"error parsing multiaddress from %s: %s",
				p.peerstorePath,
				err,
			)
		}
		addrs = append(addrs, addr)
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("reading %s: %s", p.peerstorePath, err)
	}
	return addrs
}

// SavePeerstore stores a slice of multiaddresses in the peerstore file, one
// per line.
func (p *PeerManager) SavePeerstore(pinfos []peer.AddrInfo) error {
	if p.peerstorePath == "" {
		return nil
	}

	p.peerstoreLock.Lock()
	defer p.peerstoreLock.Unlock()

	f, err := os.Create(p.peerstorePath)
	if err != nil {
		log.Errorf(
			"could not save peer addresses to %s: %s",
			p.peerstorePath,
			err,
		)
		return err
	}
	defer f.Close()

	for _, pinfo := range pinfos {
		if len(pinfo.Addrs) == 0 {
			log.Warn("address info does not have any multiaddresses")
			continue
		}

		addrs, err := peer.AddrInfoToP2pAddrs(&pinfo)
		if err != nil {
			log.Warn(err)
			continue
		}
		for _, a := range addrs {
			_, err = f.Write([]byte(fmt.Sprintf("%s\n", a.String())))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SavePeerstoreForPeers calls PeerInfos and then saves the peerstore
// file using the result.
func (p *PeerManager) SavePeerstoreForPeers(peers []peer.ID) error {
	return p.SavePeerstore(p.PeerInfos(peers))
}

// Bootstrap attempts to get up to "count" connected peers by trying those
// in the peerstore in priority order. It returns the list of peers it managed
// to connect to.
func (p *PeerManager) Bootstrap(count int) []peer.ID {
	knownPeers := p.host.Peerstore().PeersWithAddrs()
	toSort := &peerSort{
		pinfos: p2pstore.PeerInfos(p.host.Peerstore(), knownPeers),
		pstore: p.host.Peerstore(),
	}

	// Sort from highest to lowest priority
	sort.Sort(toSort)

	pinfos := toSort.pinfos
	lenKnown := len(pinfos)
	totalConns := 0
	var connectedPeers []peer.ID

	// keep conecting while we have peers in the store
	// and we have not reached count.
	for i := 0; i < lenKnown && totalConns < count; i++ {
		pinfo := pinfos[i]
		ctx, cancel := context.WithTimeout(p.ctx, ConnectTimeout)
		defer cancel()

		if p.host.Network().Connectedness(pinfo.ID) == network.Connected {
			// We are connected, assume success and do not try
			// to re-connect
			totalConns++
			continue
		}

		log.Debugf("connecting to %s", pinfo.ID)
		err := p.host.Connect(ctx, pinfo)
		if err != nil {
			log.Debug(err)
			err := p.SetPriority(pinfo.ID, 9999)
			if err != nil {
				log.Debug(err)
			}
			continue
		}
		log.Debugf("connected to %s", pinfo.ID)
		totalConns++
		connectedPeers = append(connectedPeers, pinfo.ID)
	}
	return connectedPeers
}

// SetPriority attaches a priority to a peer. 0 means more priority than
// 1. 1 means more priority than 2 etc.
func (p *PeerManager) SetPriority(pid peer.ID, prio int) error {
	return p.host.Peerstore().Put(pid, PriorityTag, prio)
}

// HandlePeerFound implements the Notifee interface for discovery.
func (p *PeerManager) HandlePeerFound(pr peer.AddrInfo) {
	addrs, err := peer.AddrInfoToP2pAddrs(&pr)
	if err != nil {
		log.Error(err)
		return
	}
	// actually mdns returns a single address but let's do things
	// as if there were several
	for _, a := range addrs {
		_, err = p.ImportPeer(a, true, peerstore.ConnectedAddrTTL)
		if err != nil {
			log.Error(err)
		}
	}
}

// peerSort is used to sort a slice of PinInfos given the PriorityTag in the
// peerstore, from the lowest tag value (0 is the highest priority) to the
// highest, Peers without a valid priority tag are considered as having a tag
// with value 0, so they will be among the first elements in the resulting
// slice.
type peerSort struct {
	pinfos []peer.AddrInfo
	pstore peerstore.Peerstore
}

func (ps *peerSort) Len() int {
	return len(ps.pinfos)
}

func (ps *peerSort) Less(i, j int) bool {
	pinfo1 := ps.pinfos[i]
	pinfo2 := ps.pinfos[j]

	var prio1, prio2 int

	prio1iface, err := ps.pstore.Get(pinfo1.ID, PriorityTag)
	if err == nil {
		prio1 = prio1iface.(int)
	}
	prio2iface, err := ps.pstore.Get(pinfo2.ID, PriorityTag)
	if err == nil {
		prio2 = prio2iface.(int)
	}
	return prio1 < prio2
}

func (ps *peerSort) Swap(i, j int) {
	pinfo1 := ps.pinfos[i]
	pinfo2 := ps.pinfos[j]
	ps.pinfos[i] = pinfo2
	ps.pinfos[j] = pinfo1
}
