package crypto

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/zenazn/pkcs7pad"
)

func AesCBCEncrypt(data, key []byte) (encrypted []byte, err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	aligned := pkcs7pad.Pad(data, aes.BlockSize)
	b, err := aes.NewCipher(key[:32])
	if err != nil {
		return
	}
	c := cipher.NewCBCEncrypter(b, key[32:])
	encrypted = make([]byte, len(aligned))
	c.CryptBlocks(encrypted, aligned)
	return
}

func AesCBCDecrypt(data, key []byte) (decrypted []byte, err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	b, err := aes.NewCipher(key[:32])
	if err != nil {
		return
	}
	c := cipher.NewCBCDecrypter(b, key[32:])
	decrypted = make([]byte, len(data))
	c.CryptBlocks(decrypted, data)
	return
}
