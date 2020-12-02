package serializers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Field string `json:"field"`
}

func TestJSONSerializer(t *testing.T) {
	s := testStruct{
		Field: "hello",
	}

	var s2 testStruct

	j := NewJSONSerializer()

	b, err := j.Serialize(s)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"field":"hello"}`), b)

	err = j.Deserialize(b, &s2)
	assert.Nil(t, err)

	assert.Equal(t, s.Field, s2.Field)
}
