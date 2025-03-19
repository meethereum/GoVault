package p2p

import (
	"fmt"
	"io"
)

const lenDecoderLimit = 5

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(r io.Reader, rpc any) error {
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	if m, ok := rpc.(*RPC); ok {
		m.Payload = buf[:n]
	} else {
		return fmt.Errorf("invalid message type")
	}

	return nil
}
