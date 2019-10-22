package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/glvd/cluster/version"
	"github.com/ipfs/go-datastore"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/ipfs-cluster/pstoremgr"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type Cluster struct {
	ctx    context.Context
	cancel context.CancelFunc

	id        peer.ID
	config    *Config
	datastore datastore.Datastore
}
type Options func(cluster *Cluster)

func DataStoreOption(datastore datastore.Datastore) Options {
	return func(cluster *Cluster) {
		cluster.datastore = datastore
	}
}

// NewCluster builds a new IPFS Cluster peer. It initializes a LibP2P host,
// creates and RPC Server and client and sets up all components.
//
// The new cluster peer may still be performing initialization tasks when
// this call returns (consensus may still be bootstrapping). Use Cluster.Ready()
// if you need to wait until the peer is fully up.
func NewCluster(
	ctx context.Context,
	host host.Host,
	//dht *dht.IpfsDHT,
	cfg *Config,
	//datastore datastore.Datastore,
	//consensus Consensus,
	//apis []API,
	//	api *httpapi.HttpApi,
	//	tracker PinTracker,
	//	monitor PeerMonitor,
	//	allocator PinAllocator,
	//	informer Informer,
	//	tracer Tracer,
	options ...Options) (*Cluster, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	if host == nil {
		return nil, errors.New("cluster host is nil")
	}

	ctx, cancel := context.WithCancel(ctx)

	listenAddrs := ""
	for _, addr := range host.Addrs() {
		listenAddrs += fmt.Sprintf("        %s/p2p/%s\n", addr, host.ID().Pretty())
	}

	log.Infof("IPFS Cluster v%s listening on:\n%s\n", version.Version, listenAddrs)

	peerManager := pstoremgr.New(ctx, host, cfg.GetPeerstorePath())

	var mdns discovery.Service
	if cfg.MDNSInterval > 0 {
		mdns, err := discovery.NewMdnsService(ctx, host, cfg.MDNSInterval, mdnsServiceTag)
		if err != nil {
			cancel()
			return nil, err
		}
		mdns.RegisterNotifee(peerManager)
	}

	c := &Cluster{
		ctx:         ctx,
		cancel:      cancel,
		id:          host.ID(),
		config:      cfg,
		host:        host,
		dht:         dht,
		discovery:   mdns,
		datastore:   datastore,
		consensus:   consensus,
		apis:        apis,
		ipfs:        ipfs,
		tracker:     tracker,
		monitor:     monitor,
		allocator:   allocator,
		informer:    informer,
		tracer:      tracer,
		peerManager: peerManager,
		shutdownB:   false,
		removed:     false,
		doneCh:      make(chan struct{}),
		readyCh:     make(chan struct{}),
		readyB:      false,
	}

	// Import known cluster peers from peerstore file. Set
	// a non permanent TTL.
	c.peerManager.ImportPeersFromPeerstore(false, peerstore.AddressTTL)
	// Attempt to connect to some peers (up to bootstrapCount)
	connectedPeers := c.peerManager.Bootstrap(bootstrapCount)
	// We cannot warn when count is low as this as this is normal if going
	// to Join() later.
	logger.Debugf("bootstrap count %d", len(connectedPeers))
	// Log a ping metric for every connected peer. This will make them
	// visible as peers without having to wait for them to send one.
	for _, p := range connectedPeers {
		if err := c.logPingMetric(ctx, p); err != nil {
			logger.Warning(err)
		}
	}

	// Bootstrap the DHT now that we possibly have some connections
	c.dht.Bootstrap(c.ctx)

	// After setupRPC components can do their tasks with a fully operative
	// routed libp2p host with some connections and a working DHT (hopefully).
	err = c.setupRPC()
	if err != nil {
		c.Shutdown(ctx)
		return nil, err
	}
	c.setupRPCClients()

	// Note: It is very important to first call Add() once in a non-racy
	// place
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.ready(ReadyTimeout)
		c.run()
	}()

	return c, nil
}
