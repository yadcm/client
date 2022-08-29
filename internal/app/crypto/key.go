package crypto

import "lukechampine.com/blake3"

func KeyFromPassword(password string) []byte {
	var hash []byte
	h := blake3.New(48, nil)
	h.Write([]byte(password))
	hash = h.Sum(nil)
	h.Reset()
	return hash
}
