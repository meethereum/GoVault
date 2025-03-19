package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents a remote node over an established TCP connection.
type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

// NewTCPPeer creates a new TCPPeer instance.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close closes the peer connection.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// TCPTransportOpts defines configurable options for TCPTransport.
type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       DefaultDecoder
	OnPeer        func(Peer) error
}

// TCPTransport represents a TCP-based transport layer.
type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandShakeFunc
	Decoder       DefaultDecoder
	rpcch         chan RPC
	OnPeer        func(Peer) error
}

// NewTCPTransport creates a new TCPTransport instance with given options.
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		listenAddress: opts.ListenAddr,
		shakeHands:    opts.HandShakeFunc,
		Decoder:       opts.Decoder,
		rpcch:         make(chan RPC),
		OnPeer:        opts.OnPeer,
	}
}

// Consume returns a read-only channel for incoming messages.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// ListenAndAccept starts listening for incoming TCP connections.
func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	t.listener = ln
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s", conn.RemoteAddr())
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	for {
		rpc := RPC{}

		err = t.Decoder.Decode(conn, &rpc)
		if err == net.ErrClosed {
			return
		}
		if err != nil {
			fmt.Printf("TCP read error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("rpc: %+v\n", rpc)
		t.rpcch <- rpc
	}
}
