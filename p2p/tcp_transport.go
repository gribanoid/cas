package p2p

import (
	"errors"
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

func (p *tcpPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder
	OnPeer func(peer Peer) error
}

type tcpTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcCh    chan RPC

	mu   sync.RWMutex
	peer map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) Transport {
	return &tcpTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

func (t *tcpTransport) Consume() <-chan RPC {
	return t.rpcCh
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
			fmt.Printf("TCP accept: %v\n", err)
		}

		fmt.Printf("new incomig connection: %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *tcpTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %v\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		if err = t.Decoder.Decode(conn, &rpc); err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Printf("TCP read error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
	}
}
