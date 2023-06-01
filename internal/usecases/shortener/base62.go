package shortener

import (
	"fmt"
	"math/big"
	"strings"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// encodeUUID implements UUID encoding (RFC-4122) in base62 string.
func encodeUUID(uuid [16]byte) string {
	num := big.NewInt(0).SetBytes(uuid[:])
	base62 := big.NewInt(62)
	remainder := big.NewInt(0)

	var buf [22]byte
	n := len(buf)

	for num.Cmp(big.NewInt(0)) > 0 {
		num.DivMod(num, base62, remainder)
		n--
		buf[n] = base62Chars[remainder.Int64()]
	}
	return string(buf[n:])
}

// decodeUUID implements decoding base62 string to UUID (RFC-4122).
func decodeUUID(s string) ([]byte, error) {
	num := big.NewInt(0)
	base62 := big.NewInt(62)

	for _, char := range s {
		index := strings.IndexByte(base62Chars, byte(char))
		if index == -1 {
			return nil, fmt.Errorf("invalid base62 character: %c", char)
		}
		num.Mul(num, base62)
		num.Add(num, big.NewInt(int64(index)))
	}

	uuidBytes := num.Bytes()

	if len(uuidBytes) < 16 {
		return nil, fmt.Errorf("invalid UUID length: %d", len(uuidBytes))
	}

	return uuidBytes, nil
}
