package help

import (
	"context"
	"errors"
	"io/ioutil"
	"sync"
)

// Manager represents an ipfs-cluster configuration which bundles
// different ComponentConfigs object together.
// Use RegisterComponent() to add a component configurations to the
// object. Once registered, configurations will be parsed from the
// central configuration file when doing LoadJSON(), and saved to it
// when doing SaveJSON().
type Manager struct {
	ctx    context.Context
	cancel func()
	wg     sync.WaitGroup
	path   string

	// The Cluster configuration has a top-level
	// special section.
	//clusterConfig ComponentConfig

	// Holds configuration objects for components.
	//sections map[SectionType]Section

	// store originally parsed jsonConfig
	//jsonCfg *jsonConfig
	// stores original source if any
	//Source string

	//sourceRedirs int // used avoid recursive source load

	// map of components which has empty configuration
	// in JSON file
	//undefinedComps map[SectionType]map[string]bool

	// if a config has been loaded from disk, track the path
	// so it can be saved to the same place.
	//path    string
	saveMux sync.Mutex
}

// NewManager returns a correctly initialized Manager
// which is ready to accept component configurations.
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:    ctx,
		cancel: cancel,
		//undefinedComps: make(map[SectionType]map[string]bool),
		//sections:       make(map[SectionType]Section),
	}

}

// SaveJSON saves the JSON representation of the Config to
// the given path.
func (m *Manager) SaveJSON(path string) error {
	m.saveMux.Lock()
	defer m.saveMux.Unlock()

	log.Info("Saving configuration")

	if path != "" {
		m.path = path
	}

	bs, err := m.ToJSON()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(m.path, bs, 0600)
	// ToJSON provides a JSON representation of the configuration by
	// generating JSON for all componenents registered.
}
func (m *Manager) ToJSON() ([]byte, error) {
	//dir := filepath.Dir(m.path)

	//err := m.Validate()
	//if err != nil {
	//	return nil, err
	//}
	//
	//if m.Source != "" {
	//	return DefaultJSONMarshal(&jsonConfig{Source: m.Source})
	//}
	//
	//jcfg := m.jsonCfg
	//if jcfg == nil {
	//	jcfg = &jsonConfig{}
	//}
	//
	//if m.clusterConfig != nil {
	//	m.clusterConfig.SetBaseDir(dir)
	//	raw, err := m.clusterConfig.ToJSON()
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//	jcfg.Cluster = new(json.RawMessage)
	//	*jcfg.Cluster = raw
	//	logger.Debug("writing changes for cluster section")
	//}
	//
	//// Given a Section and a *jsonSection, it updates the
	//// component-configurations in the latter.
	//updateJSONConfigs := func(section Section, dest *jsonSection) error {
	//	for k, v := range section {
	//		v.SetBaseDir(dir)
	//		logger.Debugf("writing changes for %s section", k)
	//		j, err := v.ToJSON()
	//		if err != nil {
	//			return err
	//		}
	//		if *dest == nil {
	//			*dest = make(jsonSection)
	//		}
	//		jsonSection := *dest
	//		jsonSection[k] = new(json.RawMessage)
	//		*jsonSection[k] = j
	//	}
	//	return nil
	//}
	//
	//for _, t := range SectionTypes() {
	//	if t == Cluster {
	//		continue
	//	}
	//	jsection := jcfg.getSection(t)
	//	err := updateJSONConfigs(m.sections[t], jsection)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//return DefaultJSONMarshal(jcfg)
	return nil, errors.New("implements me")
}
