package shortener

import (
	"fmt"
	"math/big"
)

func encodeUUID(uuid [16]byte) string {
	var i big.Int
	i.SetBytes(uuid[:])
	return i.Text(62)
}

func decodeUUID(s string) ([]byte, error) {
	var i big.Int
	_, ok := i.SetString(s, 62)
	if !ok {
		return nil, fmt.Errorf("cannot parse base62: %q", s)
	}
	var uuid []byte
	copy(uuid, i.Bytes())

	if len(i.Bytes()) < 16 {
		return nil, fmt.Errorf("invalid UUID length: %d", len(i.Bytes()))
	}
	return i.Bytes(), nil
}
