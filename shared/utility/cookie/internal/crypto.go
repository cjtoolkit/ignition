package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
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

func Sign(keyStr string, message []byte) []byte {
	if message == nil {
		return nil
	}

	key, _ := hex.DecodeString(keyStr)
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	messageMac := mac.Sum(nil)

	return append(messageMac, message...)
}

func Check(keyStr string, message []byte) []byte {
	if message == nil {
		return nil
	}
	if len(message) < sha512.Size {
		return nil
	}

	messageMac := message[:sha512.Size]
	message = message[sha512.Size:]

	key, _ := hex.DecodeString(keyStr)
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(messageMac, expectedMac) {
		return nil
	}

	return message
}
