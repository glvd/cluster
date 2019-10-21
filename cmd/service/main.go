package main

import (
	"github.com/blang/semver"
	"github.com/godcong/go-trait"
)

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

func init() {
	// Set build information.
	if build, err := semver.NewBuildVersion(commit); err == nil {
		version.Version.Build = []string{"git" + build}
	}

	// We try guessing user's home from the HOME variable. This
	// allows HOME hacks for things like Snapcraft builds. HOME
	// should be set in all UNIX by the OS. Alternatively, we fall back to
	// usr.HomeDir (which should work on Windows etc.).
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			panic(fmt.Sprintf("cannot get current user: %s", err))
		}
		home = usr.HomeDir
	}

	DefaultPath = filepath.Join(home, DefaultFolder)
}
func main() {

}
