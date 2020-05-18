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
	if value == nil {
		return nil
	}

	key, _ := hex.DecodeString(keyStr)

	bReader := bytes.NewReader(value)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}

	stream := cipher.NewOFB(block, iv)

	var out bytes.Buffer
	(&out).Write(iv)

	writer := &cipher.StreamWriter{S: stream, W: &out}
	if _, err := io.Copy(writer, bReader); err != nil {
		return nil
	}

	return out.Bytes()
}

func Decrypt(keyStr string, value []byte) []byte {
	if value == nil {
		return nil
	}

	if len(value) < aes.BlockSize {
		return nil
	}

	key, _ := hex.DecodeString(keyStr)

	iv := value[:aes.BlockSize]
	value = value[aes.BlockSize:]

	bReader := bytes.NewReader(value)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	stream := cipher.NewOFB(block, iv)

	var out bytes.Buffer

	reader := &cipher.StreamReader{S: stream, R: bReader}
	if _, err := io.Copy(&out, reader); err != nil {
		return nil
	}

	return out.Bytes()
}
