package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

	cipherText := make([]byte, aes.BlockSize+len(value))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, iv)

	var out bytes.Buffer

	writer := &cipher.StreamWriter{S: stream, W: &out}
	if _, err := io.Copy(writer, bReader); err != nil {
		panic(err)
	}

	return append(iv, out.Bytes()...)
}

func Decrypt(keyStr string, value []byte) []byte {
	key, _ := hex.DecodeString(keyStr)

	if len(value) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv := value[:aes.BlockSize]
	value = value[aes.BlockSize:]

	bReader := bytes.NewReader(value)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, iv)

	var out bytes.Buffer

	reader := &cipher.StreamReader{S: stream, R: bReader}
	if _, err := io.Copy(&out, reader); err != nil {
		panic(err)
	}

	return out.Bytes()
}
