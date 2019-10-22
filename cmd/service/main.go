package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/glvd/cluster/version"
	"github.com/godcong/go-trait"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli"
)

// ProgramName of this application
const programName = `node_cluster`

const (
	stateCleanupPrompt           = "The peer state will be removed.  Existing pins may be lost."
	configurationOverwritePrompt = "The configuration file will be overwritten."
)

var commit string

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

func checkErr(doing string, err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			doing = fmt.Sprintf(doing, args...)
		}
		log.Errorf("error %s: %s\n", doing, err)
		err = locker.tryUnlock()
		if err != nil {
			log.Errorf("error releasing execution lock: %s\n", err)
		}
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = programName
	app.Usage = "Node Cluster"
	app.Version = version.Version.String()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  DefaultPath,
			Usage:  "path to the configuration and data `FOLDER`",
			EnvVar: "CLUSTER_PATH",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "forcefully proceed with some actions. i.e. overwriting configuration",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable full debug logging (very verbose)",
		},
	}

	app.Before = func(c *cli.Context) error {
		absPath, err := filepath.Abs(c.String("config"))
		if err != nil {
			return err
		}

		configPath = filepath.Join(absPath, DefaultConfigFile)
		identityPath = filepath.Join(absPath, DefaultIdentityFile)

		locker = &lock{path: absPath}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Creates a configuration and generates an identity",
			Description: fmt.Sprintf(`
This command will initialize a new %s configuration file and, if it
does already exist, generate a new %s for %s.

If the optional [source-url] is given, the generated configuration file
will refer to it. The source configuration will be fetched from its source
URL during the launch of the daemon. If not, a default standard configuration
file will be created.

In the latter case, a cluster secret will be generated as required
by %s. Alternatively, this secret can be manually
provided with --custom-secret (in which case it will be prompted), or
by setting the CLUSTER_SECRET environment variable.

The --consensus flag allows to select an alternative consensus components for
in the newly-generated configuration.

Note that the --force flag allows to overwrite an existing
configuration with default values. To generate a new identity, please
remove the %s file first and clean any Raft state.

By default, an empty peerstore file will be created too. Initial contents can
be provided with the --peers flag. Depending on the chosen consensus, the
"trusted_peers" list in the "crdt" configuration section and the
"init_peerset" list in the "raft" configuration section will be prefilled to
the peer IDs in the given multiaddresses.
`,

				DefaultConfigFile,
				DefaultIdentityFile,
				programName,
				programName,
				DefaultIdentityFile,
			),
			ArgsUsage: "[http-source-url]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "custom-secret, s",
					Usage: "prompt for the cluster secret (when no source specified)",
				},
				cli.StringFlag{
					Name:  "peers",
					Usage: "comma-separated list of multiaddresses to init with (see help)",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "overwrite configuration without prompting",
				},
			},
			Action: func(c *cli.Context) error {
				consensus := c.String("consensus")
				if consensus != "raft" && consensus != "crdt" {
					checkErr("choosing consensus", errors.New("flag value must be set to 'raft' or 'crdt'"))
				}

				cfgHelper := cmdutils.NewConfigHelper(configPath, identityPath, consensus)
				defer cfgHelper.Manager().Shutdown() // wait for saves

				configExists := false
				if _, err := os.Stat(configPath); !os.IsNotExist(err) {
					configExists = true
				}

				identityExists := false
				if _, err := os.Stat(identityPath); !os.IsNotExist(err) {
					identityExists = true
				}

				if configExists || identityExists {
					// cluster might be running
					// acquire lock for config folder
					locker.lock()
					defer locker.tryUnlock()
				}

				if configExists {
					confirm := fmt.Sprintf(
						"%s Continue? [Y/n]:",
						configurationOverwritePrompt,
					)

					// --force allows override of the prompt
					if !c.Bool("force") {
						if !yesNoPrompt(confirm) {
							return nil
						}
					}
				}

				// Set url. If exists, it will be the only thing saved.
				cfgHelper.Manager().Source = c.Args().First()

				// Generate defaults for all registered components
				err := cfgHelper.Manager().Default()
				checkErr("generating default configuration", err)
				err = cfgHelper.Manager().ApplyEnvVars()
				checkErr("applying environment variables to configuration", err)

				//userSecret, userSecretDefined := userProvidedSecret(c.Bool("custom-secret") && !c.Args().Present())
				// Set user secret
				//if userSecretDefined {
				//	cfgHelper.Configs().Cluster.Secret = userSecret
				//}

				peersOpt := c.String("peers")
				var multiAddrs []multiaddr.Multiaddr
				if peersOpt != "" {
					addrs := strings.Split(peersOpt, ",")

					for _, addr := range addrs {
						addr = strings.TrimSpace(addr)
						multiAddr, err := multiaddr.NewMultiaddr(addr)
						checkErr("parsing peer multiaddress: "+addr, err)
						multiAddrs = append(multiAddrs, multiAddr)
					}

					peers := ipfscluster.PeersFromMultiaddrs(multiAddrs)
					cfgHelper.Configs().Crdt.TrustAll = false
					cfgHelper.Configs().Crdt.TrustedPeers = peers
					cfgHelper.Configs().Raft.InitPeerset = peers
				}

				// Save config. Creates the folder.
				// Sets BaseDir in components.
				checkErr("saving default configuration", cfgHelper.SaveConfigToDisk())
				log.Errorf("configuration written to %s.\n", configPath)

				if !identityExists {
					ident := cfgHelper.Identity()
					err := ident.Default()
					checkErr("generating an identity", err)

					err = ident.ApplyEnvVars()
					checkErr("applying environment variables to the identity", err)

					err = cfgHelper.SaveIdentityToDisk()
					checkErr("saving "+DefaultIdentityFile, err)
					log.Errorf("new identity written to %s\n", identityPath)
				}

				// Initialize peerstore file - even if empty
				peerstorePath := cfgHelper.Configs().Cluster.GetPeerstorePath()
				peerManager := pstoremgr.New(context.Background(), nil, peerstorePath)
				addrInfos, err := peer.AddrInfosFromP2pAddrs(multiAddrs...)
				checkErr("getting AddrInfos from peer multiaddresses", err)
				err = peerManager.SavePeerstore(addrInfos)
				checkErr("saving peers to peerstore", err)
				if l := len(multiAddrs); l > 0 {
					log.Errorf("peerstore written to %s with %d entries.\n", peerstorePath, len(multiAddrs))
				} else {
					log.Errorf("new empty peerstore written to %s.\n", peerstorePath)
				}

				return nil
			},
		},
		{
			Name:  "daemon",
			Usage: "Runs the IPFS Cluster peer (default)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "upgrade, u",
					Usage: "run state migrations before starting (deprecated/unused)",
				},
				cli.BoolFlag{
					Name:  "stats",
					Usage: "enable stats collection",
				},
				cli.BoolFlag{
					Name:  "tracing",
					Usage: "enable tracing collection",
				},
				cli.BoolFlag{
					Name:  "no-trust",
					Usage: "do not trust bootstrap peers (only for \"crdt\" consensus)",
				},
			},
			Action: daemon,
		},
		{
			Name:  "state",
			Usage: "Manages the peer's consensus state (pinset)",
			Subcommands: []cli.Command{
				{
					Name:  "export",
					Usage: "save the state to a JSON file",
					Description: `
This command dumps the current cluster pinset (state) as a JSON file. The
resulting file can be used to migrate, restore or backup a Cluster peer.
By default, the state will be printed to stdout.
`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Value: "",
							Usage: "writes to an output file",
						},
					},
					Action: func(c *cli.Context) error {
						locker.lock()
						defer locker.tryUnlock()

						//mgr := getStateManager()
						//
						//var w io.WriteCloser
						//var err error
						//outputPath := c.String("file")
						//if outputPath == "" {
						//	// Output to stdout
						//	w = os.Stdout
						//} else {
						//	// Create the export file
						//	w, err = os.Create(outputPath)
						//	checkErr("creating output file", err)
						//}
						//defer w.Close()

						//checkErr("exporting state", mgr.ExportState(w))
						//log.Info("state successfully exported")
						return nil
					},
				},
				{
					Name:  "import",
					Usage: "load the state from a file produced by 'export'",
					Description: `
This command reads in an exported pinset (state) file and replaces the
existing one. This can be used, for example, to restore a Cluster peer from a
backup.

If an argument is provided, it will be treated it as the path of the file
to import. If no argument is provided, stdin will be used.
`,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "skips confirmation prompt",
						},
					},
					Action: func(c *cli.Context) error {
						locker.lock()
						defer locker.tryUnlock()

						//confirm := "The pinset (state) of this peer "
						//confirm += "will be replaced. Continue? [Y/n]:"
						//if !c.Bool("force") && !yesNoPrompt(confirm) {
						//	return nil
						//}
						//
						//mgr := getStateManager()
						//
						//// Get the importing file path
						//importFile := c.Args().First()
						//var r io.ReadCloser
						//var err error
						//if importFile == "" {
						//	r = os.Stdin
						//	fmt.Println("reading from stdin, Ctrl-D to finish")
						//} else {
						//	r, err = os.Open(importFile)
						//	checkErr("reading import file", err)
						//}
						//defer r.Close()
						//
						//checkErr("importing state", mgr.ImportState(r))
						//logger.Info("state successfully imported.  Make sure all peers have consistent states")
						return nil
					},
				},
				{
					Name:  "cleanup",
					Usage: "remove persistent data",
					Description: `
This command removes any persisted consensus data in this peer, including the
current pinset (state). The next start of the peer will be like the first start
to all effects. Peers may need to bootstrap and sync from scratch after this.
`,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "skip confirmation prompt",
						},
					},
					Action: func(c *cli.Context) error {
						locker.lock()
						defer locker.tryUnlock()

						//confirm := fmt.Sprintf(
						//	"%s Continue? [Y/n]:",
						//	stateCleanupPrompt,
						//)
						//if !c.Bool("force") && !yesNoPrompt(confirm) {
						//	return nil
						//}
						//
						//mgr := getStateManager()
						//checkErr("cleaning state", mgr.Clean())
						//logger.Info("data correctly cleaned up")
						return nil
					},
				},
			},
		},
		{
			Name:  "version",
			Usage: "Prints the ipfs-cluster version",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s\n", version.Version)
				return nil
			},
		},
	}

	app.Action = run

	app.Run(os.Args)
}

// run daemon() by default, or error.
func run(c *cli.Context) error {
	cli.ShowAppHelp(c)
	os.Exit(1)
	return nil
}

// Lifted from go-ipfs/cmd/ipfs/daemon.go
func yesNoPrompt(prompt string) bool {
	var s string
	for i := 0; i < 3; i++ {
		fmt.Printf("%s ", prompt)
		fmt.Scanf("%s", &s)
		switch s {
		case "n", "N":
			return false
		}
	}
	return true
}
