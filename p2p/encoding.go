package p2p

import (
	"io"
)

const lenDecoderLimit = 5

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	rpc.Payload = buf[:n] 

	return nil
}
