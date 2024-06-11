package p2p

import "io"

type Decoder interface {
	Decode(dst any, src io.Reader) error
}

type Encoder interface{}
