package cluster

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

// connMgrConfigJSON configures the libp2p host connection manager.
type connMgrConfigJSON struct {
	HighWater   int    `json:"high_water"`
	LowWater    int    `json:"low_water"`
	GracePeriod string `json:"grace_period"`
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

func (cfg *Config) ConfigKey() string {
	panic("implement me")
}

func (cfg *Config) ToJSON() ([]byte, error) {
	panic("implement me")
}

func (cfg *Config) Default() error {
	panic("implement me")
}

func (cfg *Config) ApplyEnvVars() error {
	panic("implement me")
}

func (cfg *Config) SetBaseDir(string) {
	panic("implement me")
}

func (cfg *Config) SaveCh() <-chan struct{} {
	panic("implement me")
}

// configJSON represents a Cluster configuration as it will look when it is
// saved using JSON. Most configuration keys are converted into simple types
// like strings, and key names aim to be self-explanatory for the user.
type configJSON struct {
	ID                   string             `json:"id,omitempty"`
	PeerName             string             `json:"peername"`
	PrivateKey           string             `json:"private_key,omitempty"`
	Secret               string             `json:"secret"`
	LeaveOnShutdown      bool               `json:"leave_on_shutdown"`
	ListenMultiAddress   string             `json:"listen_multiaddress"`
	ConnectionManager    *connMgrConfigJSON `json:"connection_manager"`
	StateSyncInterval    string             `json:"state_sync_interval"`
	IPFSSyncInterval     string             `json:"ipfs_sync_interval"`
	PinRecoverInterval   string             `json:"pin_recover_interval"`
	ReplicationFactorMin int                `json:"replication_factor_min"`
	ReplicationFactorMax int                `json:"replication_factor_max"`
	MonitorPingInterval  string             `json:"monitor_ping_interval"`
	PeerWatchInterval    string             `json:"peer_watch_interval"`
	MDNSInterval         string             `json:"mdns_interval"`
	DisableRepinning     bool               `json:"disable_repinning"`
	FollowerMode         bool               `json:"follower_mode,omitempty"`
	PeerStoreFile        string             `json:"peerstore_file,omitempty"`
}

var ConfigPath string

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

// Validate will check that the values of this config
// seem to be working ones.
func (cfg *Config) Validate() error {
	if cfg.ListenAddr == nil {
		return errors.New("cluster.listen_multiaddress is undefined")
	}

	if cfg.ConnMgr.LowWater <= 0 {
		return errors.New("cluster.connection_manager.low_water is invalid")
	}

	if cfg.ConnMgr.HighWater <= 0 {
		return errors.New("cluster.connection_manager.high_water is invalid")
	}

	if cfg.ConnMgr.LowWater > cfg.ConnMgr.HighWater {
		return errors.New("cluster.connection_manager.low_water is greater than high_water")
	}

	if cfg.ConnMgr.GracePeriod == 0 {
		return errors.New("cluster.connection_manager.grace_period is invalid")
	}

	if cfg.StateSyncInterval <= 0 {
		return errors.New("cluster.state_sync_interval is invalid")
	}

	if cfg.IPFSSyncInterval <= 0 {
		return errors.New("cluster.ipfs_sync_interval is invalid")
	}

	if cfg.PinRecoverInterval <= 0 {
		return errors.New("cluster.pin_recover_interval is invalid")
	}

	if cfg.MonitorPingInterval <= 0 {
		return errors.New("cluster.monitoring_interval is invalid")
	}

	if cfg.PeerWatchInterval <= 0 {
		return errors.New("cluster.peer_watch_interval is invalid")
	}

	if err := isReplicationFactorValid(cfg.ReplicationFactorMin, cfg.ReplicationFactorMax); err != nil {
		return err
	}

	return isRPCPolicyValid(cfg.RPCPolicy)
}
func isReplicationFactorValid(rplMin, rplMax int) error {
	// check Max and Min are correct
	if rplMin == 0 || rplMax == 0 {
		return errors.New("cluster.replication_factor_min and max must be set")
	}

	if rplMin > rplMax {
		return errors.New("cluster.replication_factor_min is larger than max")
	}

	if rplMin < -1 {
		return errors.New("cluster.replication_factor_min is wrong")
	}

	if rplMax < -1 {
		return errors.New("cluster.replication_factor_max is wrong")
	}

	if (rplMin == -1 && rplMax != -1) || (rplMin != -1 && rplMax == -1) {
		return errors.New("cluster.replication_factor_min and max must be -1 when one of them is")
	}
	return nil
}
func isRPCPolicyValid(p map[string]RPCEndpointType) error {
	rpcComponents := []interface{}{
		&RPCAPI{},
		//&PinTrackerRPCAPI{},
		//&IPFSConnectorRPCAPI{},
		//&ConsensusRPCAPI{},
		//&PeerMonitorRPCAPI{},
	}

	total := 0
	for _, c := range rpcComponents {
		t := reflect.TypeOf(c)
		for i := 0; i < t.NumMethod(); i++ {
			total++
			method := t.Method(i)
			name := fmt.Sprintf("%s.%s", RPCServiceID(c), method.Name)
			_, ok := p[name]
			if !ok {
				return fmt.Errorf("RPCPolicy is missing the %s method", name)
			}
		}
	}
	if len(p) != total {
		log.Warn("defined RPC policy has more entries than needed")
	}
	return nil
}

func (cfg *Config) LoadJSON(raw []byte) error {
	jcfg := &configJSON{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		log.Error("Error unmarshaling cluster config")
		return err
	}

	return cfg.applyConfigJSON(jcfg)

}

func (cfg *Config) applyConfigJSON(jcfg *configJSON) error {
	SetIfNotDefault(jcfg.PeerStoreFile, &cfg.PeerstoreFile)

	SetIfNotDefault(jcfg.PeerName, &cfg.Peername)

	clusterSecret, err := DecodeClusterSecret(jcfg.Secret)
	if err != nil {
		err = fmt.Errorf("error loading cluster secret from config: %s", err)
		return err
	}
	cfg.Secret = clusterSecret

	clusterAddr, err := multiaddr.NewMultiaddr(jcfg.ListenMultiAddress)
	if err != nil {
		err = fmt.Errorf("error parsing cluster_listen_multiaddress: %s", err)
		return err
	}
	cfg.ListenAddr = clusterAddr

	if conman := jcfg.ConnectionManager; conman != nil {
		cfg.ConnMgr = &ConnMgrConfig{
			HighWater: jcfg.ConnectionManager.HighWater,
			LowWater:  jcfg.ConnectionManager.LowWater,
		}
		err = ParseDurations("cluster",
			&DurationOpt{Duration: jcfg.ConnectionManager.GracePeriod, Dst: &cfg.ConnMgr.GracePeriod, Name: "connection_manager.grace_period"},
		)
		if err != nil {
			return err
		}
	}

	rplMin := jcfg.ReplicationFactorMin
	rplMax := jcfg.ReplicationFactorMax
	SetIfNotDefault(rplMin, &cfg.ReplicationFactorMin)
	SetIfNotDefault(rplMax, &cfg.ReplicationFactorMax)

	err = ParseDurations("cluster",
		&DurationOpt{Duration: jcfg.StateSyncInterval, Dst: &cfg.StateSyncInterval, Name: "state_sync_interval"},
		&DurationOpt{Duration: jcfg.IPFSSyncInterval, Dst: &cfg.IPFSSyncInterval, Name: "ipfs_sync_interval"},
		&DurationOpt{Duration: jcfg.PinRecoverInterval, Dst: &cfg.PinRecoverInterval, Name: "pin_recover_interval"},
		&DurationOpt{Duration: jcfg.MonitorPingInterval, Dst: &cfg.MonitorPingInterval, Name: "monitor_ping_interval"},
		&DurationOpt{Duration: jcfg.PeerWatchInterval, Dst: &cfg.PeerWatchInterval, Name: "peer_watch_interval"},
		&DurationOpt{Duration: jcfg.MDNSInterval, Dst: &cfg.MDNSInterval, Name: "mdns_interval"},
	)
	if err != nil {
		return err
	}

	cfg.LeaveOnShutdown = jcfg.LeaveOnShutdown
	cfg.DisableRepinning = jcfg.DisableRepinning
	cfg.FollowerMode = jcfg.FollowerMode

	return nil
	//return cfg.Validate()
}

// DecodeClusterSecret parses a hex-encoded string, checks that it is exactly
// 32 bytes long and returns its value as a byte-slice.x
func DecodeClusterSecret(hexSecret string) ([]byte, error) {
	secret, err := hex.DecodeString(hexSecret)
	if err != nil {
		return nil, err
	}
	switch secretLen := len(secret); secretLen {
	case 0:
		log.Warn("Cluster secret is empty, cluster will start on unprotected network.")
		return nil, nil
	case 32:
		return secret, nil
	default:
		return nil, fmt.Errorf("input secret is %d bytes, cluster secret should be 32", secretLen)
	}
}

// GetPeerstorePath returns the full path of the
// PeerstoreFile, obtained by concatenating that value
// with BaseDir of the configuration, if set.
// An empty string is returned when BaseDir is not set.
func (cfg *Config) GetPeerstorePath() string {
	if ConfigPath == "" {
		return ""
	}

	filename := DefaultPeerstoreFile
	if cfg.PeerstoreFile != "" {
		filename = cfg.PeerstoreFile
	}

	return filepath.Join(ConfigPath, filename)
}

// DefaultJSONMarshal produces pretty JSON with 2-space indentation
func DefaultJSONMarshal(v interface{}) ([]byte, error) {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// SetIfNotDefault sets dest to the value of src if src is not the default
// value of the type.
// dest must be a pointer.
func SetIfNotDefault(src interface{}, dest interface{}) {
	switch src.(type) {
	case time.Duration:
		t := src.(time.Duration)
		if t != 0 {
			*dest.(*time.Duration) = t
		}
	case string:
		str := src.(string)
		if str != "" {
			*dest.(*string) = str
		}
	case uint64:
		n := src.(uint64)
		if n != 0 {
			*dest.(*uint64) = n
		}
	case int:
		n := src.(int)
		if n != 0 {
			*dest.(*int) = n
		}
	case bool:
		b := src.(bool)
		if b {
			*dest.(*bool) = b
		}
	}
}

// DurationOpt provides a datatype to use with ParseDurations
type DurationOpt struct {
	// The duration we need to parse
	Duration string
	// Where to store the result
	Dst *time.Duration
	// A variable name associated to it for helpful errors.
	Name string
}

// ParseDurations takes a time.Duration src and saves it to the given dst.
func ParseDurations(component string, args ...*DurationOpt) error {
	for _, arg := range args {
		if arg.Duration == "" {
			// don't do anything. Let the destination field
			// stay at its default.
			continue
		}
		t, err := time.ParseDuration(arg.Duration)
		if err != nil {
			return fmt.Errorf(
				"error parsing %s.%s: %s",
				component,
				arg.Name,
				err,
			)
		}
		*arg.Dst = t
	}
	return nil
}
