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

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type tcpTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu   sync.RWMutex
	peer map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) Transport {
	return &tcpTransport{
		TCPTransportOpts: opts,
	}
}

func (t *tcpTransport) ListenAndAccept() error {
	var err error

	if t.listener, err = net.Listen("tcp", t.ListenAddr); err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *tcpTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept: %s", err)
		}

		fmt.Printf("new incomig connection: %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *tcpTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Println("TCP handshake error:", err.Error())
		return
	}

	// Read loop
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP decoder.Decode: %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", msg)
	}
}
