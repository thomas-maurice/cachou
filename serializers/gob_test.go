package serializers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGobSerializer(t *testing.T) {
	s := testStruct{
		Field: "hello",
	}

	var s2 testStruct

	j := NewGobSerializer()

	b, err := j.Serialize(s)
	assert.Nil(t, err)

	err = j.Deserialize(b, &s2)
	assert.Nil(t, err)

	assert.Equal(t, s.Field, s2.Field)
}
