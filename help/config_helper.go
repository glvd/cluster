package help

import (
	"os"
	"path/filepath"
)

type ConfigHelper struct {
	configPath   string
	identityPath string
}

// NewConfigHelper creates a config helper given the paths to the
// configuration and identity files.
func NewConfigHelper(configPath, identityPath string) *ConfigHelper {
	ch := &ConfigHelper{
		configPath:   configPath,
		identityPath: identityPath,
	}
	//ch.init()
	return ch
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
	//return ch.manager.SaveJSON(ch.configPath)
	return nil
}
