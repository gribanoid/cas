package p2p

type HandshakeFunc func(peer Peer) error

func NOPHandshakeFunc(_ Peer) error { return nil }
