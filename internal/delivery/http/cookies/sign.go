package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("cookie not found")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func ReadSigned(r *http.Request, name string, secretKey string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", ErrNotFound
	}

	encodedValue, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	signedValue := string(encodedValue)

	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	sign := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSign := mac.Sum(nil)
	if !hmac.Equal([]byte(sign), expectedSign) {
		return "", ErrInvalidValue
	}

	return value, nil
}

func WriteSigned(w http.ResponseWriter, cookie http.Cookie, secretKey string) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	sign := mac.Sum(nil)
	cookie.Value = string(sign) + cookie.Value
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))
	cookie.Path = "/"
	http.SetCookie(w, &cookie)
}
