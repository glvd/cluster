package cluster

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Identity represents identity of a cluster peer for communication,
// including the Consensus component.
type Identity struct {
	ID         peer.ID
	PrivateKey crypto.PrivKey
}
