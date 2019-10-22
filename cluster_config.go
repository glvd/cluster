package cluster

import (
	"os"
	"time"

	pnet "github.com/libp2p/go-libp2p-pnet"
	"github.com/multiformats/go-multiaddr"
)

const configKey = "cluster"

// Configuration defaults
const (
	DefaultListenAddr          = "/ip4/0.0.0.0/tcp/9096"
	DefaultStateSyncInterval   = 600 * time.Second
	DefaultIPFSSyncInterval    = 130 * time.Second
	DefaultPinRecoverInterval  = 1 * time.Hour
	DefaultMonitorPingInterval = 15 * time.Second
	DefaultPeerWatchInterval   = 5 * time.Second
	DefaultReplicationFactor   = -1
	//DefaultLeaveOnShutdown     = false
	//DefaultDisableRepinning    = false
	DefaultPeerstoreFile      = "peerstore"
	DefaultConnMgrHighWater   = 400
	DefaultConnMgrLowWater    = 100
	DefaultConnMgrGracePeriod = 2 * time.Minute
	//DefaultFollowerMode        = false
	DefaultMDNSInterval = 10 * time.Second
)

// ConnMgrConfig configures the libp2p host connection manager.
type ConnMgrConfig struct {
	HighWater   int
	LowWater    int
	GracePeriod time.Duration
}

type Config struct {
	//config.Saver
	//lock          sync.Mutex
	//peerstoreLock sync.Mutex

	// User-defined peername for use as human-readable identifier.
	Peername string

	// Cluster secret for private network. Peers will be in the same cluster if and
	// only if they have the same ClusterSecret. The cluster secret must be exactly
	// 64 characters and contain only hexadecimal characters (`[0-9a-f]`).
	Secret []byte

	// RPCPolicy defines access control to RPC endpoints.
	RPCPolicy map[string]RPCEndpointType

	// Leave Cluster on shutdown. Politely informs other peers
	// of the departure and removes itself from the consensus
	// peer set. The Cluster size will be reduced by one.
	LeaveOnShutdown bool

	// Listen parameters for the Cluster libp2p Host. Used by
	// the RPC and Consensus components.
	ListenAddr multiaddr.Multiaddr

	// ConnMgr holds configuration values for the connection manager for
	// the libp2p host.
	// FIXME: This only applies to ipfs-cluster-service.
	ConnMgr *ConnMgrConfig

	// Time between syncs of the consensus state to the
	// tracker state. Normally states are synced anyway, but this helps
	// when new nodes are joining the cluster. Reduce for faster
	// consistency, increase with larger states.
	StateSyncInterval time.Duration

	// Time between syncs of the local state and
	// the state of the ipfs daemon. This ensures that cluster
	// provides the right status for tracked items (for example
	// to detect that a pin has been removed. Reduce for faster
	// consistency, increase when the number of pinned items is very
	// large.
	IPFSSyncInterval time.Duration

	// Time between automatic runs of the "recover" operation
	// which will retry to pin/unpin items in error state.
	PinRecoverInterval time.Duration

	// ReplicationFactorMax indicates the target number of nodes
	// that should pin content. For exampe, a replication_factor of
	// 3 will have cluster allocate each pinned hash to 3 peers if
	// possible.
	// See also ReplicationFactorMin. A ReplicationFactorMax of -1
	// will allocate to every available node.
	ReplicationFactorMax int

	// ReplicationFactorMin indicates the minimum number of healthy
	// nodes pinning content. If the number of nodes available to pin
	// is less than this threshold, an error will be returned.
	// In the case of peer health issues, content pinned will be
	// re-allocated if the threshold is crossed.
	// For exampe, a ReplicationFactorMin of 2 will allocate at least
	// two peer to hold content, and return an error if this is not
	// possible.
	ReplicationFactorMin int

	// MonitorPingInterval is the frequency with which a cluster peer pings
	// the monitoring component. The ping metric has a TTL set to the double
	// of this value.
	MonitorPingInterval time.Duration

	// PeerWatchInterval is the frequency that we use to watch for changes
	// in the consensus peerset and save new peers to the configuration
	// file. This also affects how soon we realize that we have
	// been removed from a cluster.
	PeerWatchInterval time.Duration

	// MDNSInterval controls the time between mDNS broadcasts to the
	// network announcing the peer addresses. Set to 0 to disable
	// mDNS.
	MDNSInterval time.Duration

	// If true, DisableRepinning, ensures that no repinning happens
	// when a node goes down.
	// This is useful when doing certain types of maintainance, or simply
	// when not wanting to rely on the monitoring system which needs a revamp.
	DisableRepinning bool

	// FollowerMode disables broadcast requests from this peer
	// (sync, recover, status) and disallows pinset management
	// operations (Pin/Unpin).
	FollowerMode bool

	// Peerstore file specifies the file on which we persist the
	// libp2p host peerstore addresses. This file is regularly saved.
	PeerstoreFile string

	// Tracing flag used to skip tracing specific paths when not enabled.
	Tracing bool
}

// Default fills in all the Config fields with
// default working values. This means, it will
// generate a Secret.
func DefaultConfig() (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	addr, _ := multiaddr.NewMultiaddr(DefaultListenAddr)

	//cluster secret
	clusterSecret, err := pnet.GenerateV1Bytes()
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Peername:        hostname,
		Secret:          (*clusterSecret)[:],
		RPCPolicy:       DefaultRPCPolicy,
		LeaveOnShutdown: false,
		ListenAddr:      addr,
		ConnMgr: &ConnMgrConfig{
			HighWater:   DefaultConnMgrHighWater,
			LowWater:    DefaultConnMgrLowWater,
			GracePeriod: DefaultConnMgrGracePeriod,
		},
		StateSyncInterval:    DefaultStateSyncInterval,
		IPFSSyncInterval:     DefaultIPFSSyncInterval,
		PinRecoverInterval:   DefaultPinRecoverInterval,
		ReplicationFactorMax: DefaultReplicationFactor,
		ReplicationFactorMin: DefaultReplicationFactor,
		MonitorPingInterval:  DefaultMonitorPingInterval,
		PeerWatchInterval:    DefaultPeerWatchInterval,
		MDNSInterval:         DefaultMDNSInterval,
		DisableRepinning:     false,
		FollowerMode:         false,
		PeerstoreFile:        "",
		Tracing:              false,
	}
	return cfg, nil
}
