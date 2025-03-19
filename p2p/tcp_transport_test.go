package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: func(p Peer) error { return nil },
		Decoder:       DefaultDecoder{},
		OnPeer:        func(p Peer) error { return nil },
	}

	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.listenAddress, opts.ListenAddr)
	assert.Nil(t, tr.ListenAndAccept())
}
