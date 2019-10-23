package cluster

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glvd/cluster/version"
	"github.com/goextension/log"
	"github.com/ipfs/go-datastore"
	ocgorpc "github.com/lanzafame/go-libp2p-ocgorpc"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	rpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/multiformats/go-multiaddr"
	"go.opencensus.io/trace"
)

const (
	pingMetricName      = "ping"
	bootstrapCount      = 3
	reBootstrapInterval = 30 * time.Second
	mdnsServiceTag      = "_cluster-discovery._udp"
)

type Cluster struct {
	ctx    context.Context
	cancel context.CancelFunc

	id          peer.ID
	config      *Config
	datastore   datastore.Datastore
	host        host.Host
	discovery   discovery.Service
	peerManager *PeerManager
	rpcServer   *rpc.Server
	rpcClient   *rpc.Client
}

type Options func(cluster *Cluster)

// ReadyTimeout specifies the time before giving up
// during startup (waiting for consensus to be ready)
// It may need adjustment according to timeouts in the
// consensus layer.
var ReadyTimeout = 30 * time.Second

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

	//peerManager := pstoremgr.New(ctx, host, cfg.GetPeerstorePath())

	//var mdns discovery.Service
	//if cfg.MDNSInterval > 0 {
	//	mdns, err := discovery.NewMdnsService(ctx, host, cfg.MDNSInterval, mdnsServiceTag)
	//	if err != nil {
	//		cancel()
	//		return nil, err
	//	}
	//mdns.RegisterNotifee(peerManager)
	//}

	c := &Cluster{
		ctx:    ctx,
		cancel: cancel,
		id:     host.ID(),
		config: cfg,
		host:   host,
		//dht:         dht,
		//discovery: mdns,
		//datastore:   datastore,
		//consensus:   consensus,
		//apis:        apis,
		//ipfs:        ipfs,
		//tracker:     tracker,
		//monitor:     monitor,
		//allocator:   allocator,
		//informer:    informer,
		//tracer:      tracer,
		//peerManager: peerManager,
		//shutdownB:   false,
		//removed:     false,
		//doneCh:      make(chan struct{}),
		//readyCh:     make(chan struct{}),
		//readyB:      false,
	}

	for _, option := range options {
		option(c)
	}
	// Import known cluster peers from peerstore file. Set
	// a non permanent TTL.
	//c.peerManager.ImportPeersFromPeerstore(false, peerstore.AddressTTL)
	// Attempt to connect to some peers (up to bootstrapCount)
	//connectedPeers := c.peerManager.Bootstrap(bootstrapCount)
	// We cannot warn when count is low as this as this is normal if going
	// to Join() later.
	//log.Debugf("bootstrap count %d", len(connectedPeers))
	// Log a ping metric for every connected peer. This will make them
	// visible as peers without having to wait for them to send one.
	//for _, p := range connectedPeers {
	//	if err := c.logPingMetric(ctx, p); err != nil {
	//		log.Warning(err)
	//	}
	//}

	// Bootstrap the DHT now that we possibly have some connections
	//c.dht.Bootstrap(c.ctx)

	// After setupRPC components can do their tasks with a fully operative
	// routed libp2p host with some connections and a working DHT (hopefully).
	//err = c.setupRPC()
	//if err != nil {
	//	c.Shutdown(ctx)
	//	return nil, err
	//}
	//c.setupRPCClients()
	//
	// Note: It is very important to first call Add() once in a non-racy
	// place
	//c.wg.Add(1)
	//go func() {
	//	defer c.wg.Done()
	//	c.ready(ReadyTimeout)
	//	c.run()
	//}()
	//
	return c, nil
}

func (c *Cluster) Join(ctx context.Context, addr multiaddr.Multiaddr) error {
	_, span := trace.StartSpan(ctx, "cluster/Join")
	defer span.End()
	ctx = trace.NewContext(c.ctx, span)

	log.Debugf("Join(%s)", addr)

	// Add peer to peerstore so we can talk to it (and connect)
	pid, err := c.peerManager.ImportPeer(addr, true, peerstore.PermanentAddrTTL)
	if err != nil {
		return err
	}
	if pid == c.id {
		return nil
	}

	// Note that PeerAdd() on the remote peer will
	// figure out what our real address is (obviously not
	// ListenAddr).
	var myID api.ID
	err = c.rpcClient.CallContext(
		ctx,
		pid,
		"Cluster",
		"PeerAdd",
		c.id,
		&myID,
	)
	if err != nil {
		log.Error(err)
		return err
	}

	// Log a fake but valid metric from the peer we are
	// contacting. This will signal a CRDT component that
	// we know that peer since we have metrics for it without
	// having to wait for the next metric round.
	//if err := c.logPingMetric(ctx, pid); err != nil {
	//	log.Warn(err)
	//}

	// Broadcast our metrics to the world
	//_, err = c.sendInformerMetric(ctx)
	//if err != nil {
	//	log.Warn(err)
	//}
	//_, err = c.sendPingMetric(ctx)
	//if err != nil {
	//	log.Warning(err)
	//}

	// We need to trigger a DHT bootstrap asap for this peer to not be
	// lost if the peer it bootstrapped to goes down. We do this manually
	// by triggering 1 round of bootstrap in the background.
	// Note that our regular bootstrap process is still running in the
	// background since we created the cluster.
	//go func() {
	//	c.dht.BootstrapOnce(ctx, dht.DefaultBootstrapConfig)
	//}()

	// ConnectSwarms in the background after a while, when we have likely
	// received some metrics.
	//time.AfterFunc(c.config.MonitorPingInterval, func() {
	//	c.ipfs.ConnectSwarms(ctx)
	//})

	// wait for leader and for state to catch up
	// then sync
	//err = c.consensus.WaitForSync(ctx)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}

	//c.StateSync(ctx)

	log.Infof("%s: joined %s's cluster", c.id.Pretty(), pid.Pretty())
	return nil
}

// Done provides a way to learn if the Peer has been shutdown
// (for example, because it has been removed from the Cluster)
func (c *Cluster) Done() <-chan struct{} {
	return nil
}

func (c *Cluster) Shutdown(ctx context.Context) error {
	return nil
}

func (c *Cluster) setupRPC() error {
	rpcServer, err := newRPCServer(c)
	if err != nil {
		return err
	}
	c.rpcServer = rpcServer

	var rpcClient *rpc.Client
	if c.config.Tracing {
		csh := &ocgorpc.ClientHandler{}
		rpcClient = rpc.NewClientWithServer(
			c.host,
			version.RPCProtocol,
			rpcServer,
			rpc.WithClientStatsHandler(csh),
		)
	} else {
		rpcClient = rpc.NewClientWithServer(c.host, version.RPCProtocol, rpcServer)
	}
	c.rpcClient = rpcClient
	return nil
}
