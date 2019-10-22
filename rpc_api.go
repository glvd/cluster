package cluster

import (
	"github.com/glvd/cluster/version"
	"github.com/libp2p/go-libp2p-core/peer"
	rpc "github.com/libp2p/go-libp2p-gorpc"

	ocgorpc "github.com/lanzafame/go-libp2p-ocgorpc"
)

// RPC endpoint types w.r.t. trust level
const (
	// RPCClosed endpoints can only be called by the local cluster peer
	// on itself.
	RPCClosed RPCEndpointType = iota
	// RPCTrusted endpoints can be called by "trusted" peers.
	// It depends which peers are considered trusted. For example,
	// in "raft" mode, Cluster will all peers as "trusted". In "crdt" mode,
	// trusted peers are those specified in the configuration.
	RPCTrusted
	// RPCOpen endpoints can be called by any peer in the Cluster swarm.
	RPCOpen
)

// RPCEndpointType controls how access is granted to an RPC endpoint
type RPCEndpointType int

// A trick to find where something is used (i.e. Cluster.Pin):
// grep -R -B 3 '"Pin"' | grep -C 1 '"Cluster"'.
// This does not cover globalPinInfo*(...) broadcasts nor redirects to leader
// in Raft.

// newRPCServer returns a new RPC Server for Cluster.
func newRPCServer(c *Cluster) (*rpc.Server, error) {
	var s *rpc.Server

	authF := func(pid peer.ID, svc, method string) bool {
		endpointType, ok := c.config.RPCPolicy[svc+"."+method]
		if !ok {
			return false
		}

		switch endpointType {
		//case RPCTrusted:
		//	return c.consensus.IsTrustedPeer(c.ctx, pid)
		case RPCOpen:
			return true
		default:
			return false
		}
	}

	if c.config.Tracing {
		s = rpc.NewServer(
			c.host,
			version.RPCProtocol,
			rpc.WithServerStatsHandler(&ocgorpc.ServerHandler{}),
			rpc.WithAuthorizeFunc(authF),
		)
	} else {
		s = rpc.NewServer(c.host, version.RPCProtocol, rpc.WithAuthorizeFunc(authF))
	}

	cl := &RPCAPI{c}
	err := s.RegisterName(RPCServiceID(cl), cl)
	if err != nil {
		return nil, err
	}
	//pt := &PinTrackerRPCAPI{c.tracker}
	//err = s.RegisterName(RPCServiceID(pt), pt)
	//if err != nil {
	//	return nil, err
	//}
	//ic := &IPFSConnectorRPCAPI{c.ipfs}
	//err = s.RegisterName(RPCServiceID(ic), ic)
	//if err != nil {
	//	return nil, err
	//}
	//cons := &ConsensusRPCAPI{c.consensus}
	//err = s.RegisterName(RPCServiceID(cons), cons)
	//if err != nil {
	//	return nil, err
	//}
	//pm := &PeerMonitorRPCAPI{c.monitor}
	//err = s.RegisterName(RPCServiceID(pm), pm)
	//if err != nil {
	//	return nil, err
	//}
	return s, nil
}

// RPCServiceID returns the Service ID for the given RPCAPI object.
func RPCServiceID(rpcSvc interface{}) string {
	if v, b := rpcSvc.(IServiceName); b {
		return v.Name()
	}

	//switch rpcSvc.(type) {
	//case *RPCAPI:
	//	return "Cluster"
	//case *PinTrackerRPCAPI:
	//	return "PinTracker"
	//case *IPFSConnectorRPCAPI:
	//	return "IPFSConnector"
	//case *ConsensusRPCAPI:
	//	return "Consensus"
	//case *PeerMonitorRPCAPI:
	//	return "PeerMonitor"
	//default:
	//	return ""
	//}
	return ""
}

type IServiceName interface {
	Name() string
}

// RPCAPI is a go-libp2p-gorpc service which provides the internal peer
// API for the main cluster component.
type RPCAPI struct {
	c *Cluster
}

func (RPCAPI) Name() string {
	return "Cluster"
}
