package cache

import (
	"bytes"
	"encoding/gob"
)

func GobMarshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

func GobUnmarshal(data []byte, v interface{}) error {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(v)
	return err
}
