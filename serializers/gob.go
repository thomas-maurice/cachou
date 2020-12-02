package serializers

import (
	"bytes"
	"encoding/gob"
)

type GobSerializer struct {
}

func NewGobSerializer() *GobSerializer {
	return &GobSerializer{}
}

func (s *GobSerializer) Serialize(obj interface{}) ([]byte, error) {
	var output bytes.Buffer
	enc := gob.NewEncoder(&output)
	err := enc.Encode(obj)
	return output.Bytes(), err
}

func (s *GobSerializer) Deserialize(data []byte, obj interface{}) error {
	enc := gob.NewDecoder(bytes.NewBuffer(data))
	return enc.Decode(obj)
}