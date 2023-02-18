package shortener

import (
	"fmt"
	"math/big"
)

func uintEncode(number uint64) string {
	return big.NewInt(int64(number)).Text(62)
}

func uintDecode(s string) (uint64, error) {
	n := new(big.Int)
	_, ok := n.SetString(s, 62)
	if !ok {
		return 0, fmt.Errorf("failed decode string %s", s)
	}

	if !n.IsUint64() {
		return 0, fmt.Errorf("%s is not uint64 encoded", s)
	}

	return n.Uint64(), nil
}
