package p2p

import (
	"fmt"
	"net"
	"sync"
)

// tcpPeer represents the remote node over a TCP established connection.
type tcpPeer struct {
	conn net.Conn

	// if we dial and retrieve conn => outbound == true
	// if we accept and retrieve conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) Peer {
	return &tcpPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type tcpTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder

	mu   sync.RWMutex
	peer map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) Transport {
	return &tcpTransport{
		shakeHands: NOPHandshakeFunc, // TODO
		listenAddr: listenAddr,
	}
}

func (t *tcpTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *tcpTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept: %s", err) // TODO
		}

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *tcpTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {

	}

	// Read loop
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(msg, conn); err != nil {
			fmt.Printf("TCP decoder.Decode: %s\n", err) // TODO
			continue
		}
	}
}
