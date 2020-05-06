package embedder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"log"
)

func DecodeValue(value string) []byte {
	raw, err := base64.StdEncoding.DecodeString(value)
	errCheck(err)

	r, err := gzip.NewReader(bytes.NewReader(raw))
	errCheck(err)

	b, err := ioutil.ReadAll(r)
	errCheck(err)

	return b
}

func DecodeValueStr(value string) string { return string(DecodeValue(value)) }

func errCheck(err error) {
	if err != nil {
		log.Panic(err)
	}
}
