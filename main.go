package main

import (
	"fmt"
	"log"

	"github.com/meethereum/GoVault/p2p"
)

func OnPeer(p2p.Peer) error{
	fmt.Println("doing some logic outside peer")
	return nil
}

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: func(p p2p.Peer) error { return nil },
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(opts)

	go func() {
		for msg := range tr.Consume() {
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {} // Keeps the program running
}