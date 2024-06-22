package main

import (
	"fmt"
	"log"

	"github.com/gribanoid/cas/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("doing some logic with the peer outside of TCPTransport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Received message: %+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
