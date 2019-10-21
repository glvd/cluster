package main

import "github.com/godcong/go-trait"

var log = trait.NewZapSugar()

// Default location for the configurations and data
var (
	// DefaultFolder is the name of the cluster folder
	DefaultFolder = ".cluster"
	// DefaultPath is set on init() to $HOME/DefaultFolder
	// and holds all the cluster data
	DefaultPath string
	// The name of the configuration file inside DefaultPath
	DefaultConfigFile = "service.json"
	// The name of the identity file inside DefaultPath
	DefaultIdentityFile = "identity.json"
)

var (
	configPath   string
	identityPath string
)

func main() {

}
