package serializers

import (
	"encoding/json"
)

type JSONSerializer struct {
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (s *JSONSerializer) Serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (s *JSONSerializer) Deserialize(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}
