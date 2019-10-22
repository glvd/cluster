package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/glvd/cluster"
	"github.com/goextension/log"
	ipfscluster "github.com/ipfs/ipfs-cluster"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli"
)

func parseBootstraps(flagVal []string) (bootstraps []multiaddr.Multiaddr) {
	for _, a := range flagVal {
		bAddr, err := multiaddr.NewMultiaddr(strings.TrimSpace(a))
		checkErr("error parsing bootstrap multiaddress (%s)", err, a)
		bootstraps = append(bootstraps, bAddr)
	}
	return
}

// Runs the cluster peer
func daemon(c *cli.Context) error {
	log.Info("Initializing. For verbose output run with \"-l debug\". Please wait...")

	ctx, cancel := context.WithCancel(context.Background())
	var bootstraps []multiaddr.Multiaddr
	//if bootStr := c.String("bootstrap"); bootStr != "" {
	//	bootstraps = parseBootstraps(strings.Split(bootStr, ","))
	//}

	// Execution lock
	locker.lock()
	defer locker.tryUnlock()

	host, pubsub, dht, err := ipfscluster.NewClusterHost(ctx, nil, nil)
	checkErr("creating libp2p host", err)

	var cluster, err = createCluster(ctx, host)
	checkErr("starting cluster", err)

	// noop if no bootstraps
	// if bootstrapping fails, consensus will never be ready
	// and timeout. So this can happen in background and we
	// avoid worrying about error handling here (since Cluster
	// will realize).
	//TODO:
	go bootstrap(ctx, cluster, bootstraps)

	return handleSignals(ctx, cancel, cluster, nil, nil)
}

// createCluster creates all the necessary things to produce the cluster
// object and returns it along the datastore so the lifecycle can be handled
// (the datastore needs to be Closed after shutting down the Cluster).
func createCluster(
	ctx context.Context,
	//c *cli.Context,
	host host.Host,
	pubsub *pubsub.PubSub,
	dht *dht.IpfsDHT,
	//raftStaging bool,
) (*cluster.Cluster, error) {
	config, e := cluster.DefaultConfig()
	if e != nil {
		return nil, e
	}
	return cluster.NewCluster(
		ctx,
		host,
		config,
		//dht,
		//store,
		//cons,
		//apis,
		//connector,
		//tracker,
		//mon,
		//alloc,
		//informer,
		//tracer,
	)
}

// bootstrap will bootstrap this peer to one of the bootstrap addresses
// if there are any.
func bootstrap(ctx context.Context, cluster *cluster.Cluster, bootstraps []multiaddr.Multiaddr) {
	for _, bstrap := range bootstraps {
		log.Infof("Bootstrapping to %s", bstrap)
		err := cluster.Join(ctx, bstrap)
		if err != nil {
			log.Errorf("bootstrap to %s failed: %s", bstrap, err)
		}
	}
}

func handleSignals(
	ctx context.Context,
	cancel context.CancelFunc,
	cluster *cluster.Cluster,
	host host.Host,
	dht *dht.IpfsDHT,
) error {
	signalChan := make(chan os.Signal, 20)
	signal.Notify(
		signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	var ctrlcCount int
	for {
		select {
		case <-signalChan:
			ctrlcCount++
			handleCtrlC(ctx, cluster, ctrlcCount)
		case <-cluster.Done():
			cancel()
			dht.Close()
			host.Close()
			return nil
		}
	}
}

func handleCtrlC(ctx context.Context, cluster *cluster.Cluster, ctrlcCount int) {
	switch ctrlcCount {
	case 1:
		go func() {
			err := cluster.Shutdown(ctx)
			checkErr("shutting down cluster", err)
		}()
	case 2:
		log.Error(`


!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
Shutdown is taking too long! Press Ctrl-c again to manually kill cluster.
Note that this may corrupt the local cluster state.
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!


`)
	case 3:
		log.Error("exiting cluster NOW")
		locker.tryUnlock()
		os.Exit(-1)
	}
}
