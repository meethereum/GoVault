package p2p

import (
	//"bytes"
	"fmt"
	"net"
	"sync"
)

// a remote node over a established TCP connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn
	// if we dial to a connection to a peer it is outbound
	// if we accept it is a inbound peer
	outbound bool
}

// constructor for tcp peer
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandShakeFunc
	Decoder       DefaultDecoder
	peers         map[net.Addr]Peer
	mu            sync.RWMutex
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands: func(Peer) error {
			return nil
		},
		listenAddress: listenAddr,
		Decoder:       DefaultDecoder{},
	}
}

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
		}
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	fmt.Printf("new incoming peer : %v\n", peer)

	for {
		rpc := &RPC{}
		lenDecoderError := 0

		if err := t.Decoder.Decode(conn, rpc); err != nil {
			lenDecoderError++
			if lenDecoderError >= lenDecoderLimit {
				//drop the connection

			}
			fmt.Printf("TCP error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		fmt.Printf("message: %+v\n", rpc)
	}

}
