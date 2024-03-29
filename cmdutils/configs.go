package cmdutils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	ipfscluster "github.com/glvd/cluster"
	"github.com/glvd/cluster/api/ipfsproxy"
	"github.com/glvd/cluster/api/rest"
	"github.com/glvd/cluster/config"
	"github.com/glvd/cluster/consensus/crdt"
	"github.com/glvd/cluster/consensus/raft"
	"github.com/glvd/cluster/datastore/badger"
	"github.com/glvd/cluster/informer/disk"
	"github.com/glvd/cluster/informer/numpin"
	"github.com/glvd/cluster/ipfsconn/ipfshttp"
	"github.com/glvd/cluster/monitor/pubsubmon"
	"github.com/glvd/cluster/observations"
	"github.com/glvd/cluster/pintracker/maptracker"
	"github.com/glvd/cluster/pintracker/stateless"
)

// Configs carries config types used by a Cluster Peer.
type Configs struct {
	Cluster          *ipfscluster.Config
	Restapi          *rest.Config
	Ipfsproxy        *ipfsproxy.Config
	Ipfshttp         *ipfshttp.Config
	Raft             *raft.Config
	Crdt             *crdt.Config
	Maptracker       *maptracker.Config
	Statelesstracker *stateless.Config
	Pubsubmon        *pubsubmon.Config
	Diskinf          *disk.Config
	Numpininf        *numpin.Config
	Metrics          *observations.MetricsConfig
	Tracing          *observations.TracingConfig
	Badger           *badger.Config
}

// ConfigHelper helps managing the configuration and identity files with the
// standard set of cluster components.
type ConfigHelper struct {
	identity *config.Identity
	manager  *config.Manager
	configs  *Configs

	configPath   string
	identityPath string
	consensus    string
}

// NewConfigHelper creates a config helper given the paths to the
// configuration and identity files.
func NewConfigHelper(configPath, identityPath, consensus string) *ConfigHelper {
	ch := &ConfigHelper{
		configPath:   configPath,
		identityPath: identityPath,
		consensus:    consensus,
	}
	ch.init()
	return ch
}

// LoadConfigFromDisk parses the configuration from disk.
func (ch *ConfigHelper) LoadConfigFromDisk() error {
	return ch.manager.LoadJSONFileAndEnv(ch.configPath)
}

// LoadIdentityFromDisk parses the identity from disk.
func (ch *ConfigHelper) LoadIdentityFromDisk() error {
	// load identity with hack for 0.11.0 - identity separation.
	_, err := os.Stat(ch.identityPath)
	ident := &config.Identity{}
	// temporary hack to convert identity
	if os.IsNotExist(err) {
		clusterConfig, err := config.GetClusterConfig(ch.configPath)
		if err != nil {
			return err
		}
		err = ident.LoadJSON(clusterConfig)
		if err != nil {
			return errors.Wrap(err, "error loading identity")
		}

		err = ident.SaveJSON(ch.identityPath)
		if err != nil {
			return errors.Wrap(err, "error saving identity")
		}

		fmt.Fprintf(
			os.Stderr,
			"\nNOTICE: identity information extracted from %s and saved as %s.\n\n",
			ch.configPath,
			ch.identityPath,
		)
	} else { // leave this part when the hack is removed.
		err = ident.LoadJSONFromFile(ch.identityPath)
		if err != nil {
			return fmt.Errorf("error loading identity from %s: %s", ch.identityPath, err)
		}
	}

	err = ident.ApplyEnvVars()
	if err != nil {
		return errors.Wrap(err, "error applying environment variables to the identity")
	}
	ch.identity = ident
	return nil
}

// LoadFromDisk loads both configuration and identity from disk.
func (ch *ConfigHelper) LoadFromDisk() error {
	err := ch.LoadConfigFromDisk()
	if err != nil {
		return err
	}
	return ch.LoadIdentityFromDisk()
}

// Identity returns the Identity object. It returns an empty identity
// if not loaded yet.
func (ch *ConfigHelper) Identity() *config.Identity {
	return ch.identity
}

// Manager returns the config manager with all the
// cluster configurations registered.
func (ch *ConfigHelper) Manager() *config.Manager {
	return ch.manager
}

// Configs returns the Configs object which holds all the cluster
// configurations. Configurations are empty if they have not been loaded from
// disk.
func (ch *ConfigHelper) Configs() *Configs {
	return ch.configs
}

// GetConsensus attempts to return the configured consensus.
// If the ConfigHelper was initialized with a consensus string
// then it returns that.
//
// Otherwise it checks whether one of the consensus configurations
// has been loaded. If both or non have been loaded, it returns
// an empty string.
func (ch *ConfigHelper) GetConsensus() string {
	if ch.consensus != "" {
		return ch.consensus
	}
	crdtLoaded := ch.manager.IsLoadedFromJSON(config.Consensus, ch.configs.Crdt.ConfigKey())
	raftLoaded := ch.manager.IsLoadedFromJSON(config.Consensus, ch.configs.Raft.ConfigKey())
	if crdtLoaded == raftLoaded { //both loaded or none
		return ""
	}

	if crdtLoaded {
		return ch.configs.Crdt.ConfigKey()
	}
	return ch.configs.Raft.ConfigKey()
}

// register all current cluster components
func (ch *ConfigHelper) init() {
	man := config.NewManager()
	cfgs := &Configs{
		Cluster:          &ipfscluster.Config{},
		Restapi:          &rest.Config{},
		Ipfsproxy:        &ipfsproxy.Config{},
		Ipfshttp:         &ipfshttp.Config{},
		Raft:             &raft.Config{},
		Crdt:             &crdt.Config{},
		Maptracker:       &maptracker.Config{},
		Statelesstracker: &stateless.Config{},
		Pubsubmon:        &pubsubmon.Config{},
		Diskinf:          &disk.Config{},
		Metrics:          &observations.MetricsConfig{},
		Tracing:          &observations.TracingConfig{},
		Badger:           &badger.Config{},
	}
	man.RegisterComponent(config.Cluster, cfgs.Cluster)
	man.RegisterComponent(config.API, cfgs.Restapi)
	man.RegisterComponent(config.API, cfgs.Ipfsproxy)
	man.RegisterComponent(config.IPFSConn, cfgs.Ipfshttp)
	man.RegisterComponent(config.PinTracker, cfgs.Maptracker)
	man.RegisterComponent(config.PinTracker, cfgs.Statelesstracker)
	man.RegisterComponent(config.Monitor, cfgs.Pubsubmon)
	man.RegisterComponent(config.Informer, cfgs.Diskinf)
	man.RegisterComponent(config.Observations, cfgs.Metrics)
	man.RegisterComponent(config.Observations, cfgs.Tracing)

	switch ch.consensus {
	case cfgs.Raft.ConfigKey():
		man.RegisterComponent(config.Consensus, cfgs.Raft)
	case cfgs.Crdt.ConfigKey():
		man.RegisterComponent(config.Consensus, cfgs.Crdt)
		man.RegisterComponent(config.Datastore, cfgs.Badger)
	default:
		man.RegisterComponent(config.Consensus, cfgs.Raft)
		man.RegisterComponent(config.Consensus, cfgs.Crdt)
		man.RegisterComponent(config.Datastore, cfgs.Badger)
	}

	ch.identity = &config.Identity{}
	ch.manager = man
	ch.configs = cfgs
}

// MakeConfigFolder creates the folder to hold
// configuration and identity files.
func (ch *ConfigHelper) MakeConfigFolder() error {
	f := filepath.Dir(ch.configPath)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		err := os.MkdirAll(f, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveConfigToDisk saves the configuration file to disk.
func (ch *ConfigHelper) SaveConfigToDisk() error {
	err := ch.MakeConfigFolder()
	if err != nil {
		return err
	}
	return ch.manager.SaveJSON(ch.configPath)
}

// SaveIdentityToDisk saves the identity file to disk.
func (ch *ConfigHelper) SaveIdentityToDisk() error {
	err := ch.MakeConfigFolder()
	if err != nil {
		return err
	}
	return ch.Identity().SaveJSON(ch.identityPath)
}

// SetupTracing propagates tracingCfg.EnableTracing to all other
// configurations. Use only when identity has been loaded or generated.  The
// forceEnabled parameter allows to override the EnableTracing value.
func (ch *ConfigHelper) SetupTracing(forceEnabled bool) {
	enabled := forceEnabled || ch.configs.Tracing.EnableTracing

	ch.configs.Tracing.ClusterID = ch.Identity().ID.Pretty()
	ch.configs.Tracing.ClusterPeername = ch.configs.Cluster.Peername
	ch.configs.Tracing.EnableTracing = enabled
	ch.configs.Cluster.Tracing = enabled
	ch.configs.Raft.Tracing = enabled
	ch.configs.Crdt.Tracing = enabled
	ch.configs.Restapi.Tracing = enabled
	ch.configs.Ipfshttp.Tracing = enabled
	ch.configs.Ipfsproxy.Tracing = enabled
}
