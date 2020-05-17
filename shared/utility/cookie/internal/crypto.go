package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"io"
)

func Encrypt(keyStr string, value []byte) []byte {
	key, _ := hex.DecodeString(keyStr)

	bReader := bytes.NewReader(value)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	var out bytes.Buffer

	writer := &cipher.StreamWriter{S: stream, W: &out}
	if _, err := io.Copy(writer, bReader); err != nil {
		panic(err)
	}

	return out.Bytes()
}

func Decrypt(keyStr string, value []byte) []byte {
	key, _ := hex.DecodeString(keyStr)

	bReader := bytes.NewReader(value)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	var out bytes.Buffer

	reader := &cipher.StreamReader{S: stream, R: bReader}
	if _, err := io.Copy(&out, reader); err != nil {
		panic(err)
	}

	return out.Bytes()
}
