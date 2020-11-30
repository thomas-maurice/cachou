package cachou

import (
	"os"
	"path"
	"testing"

	"github.com/thomas-maurice/cachou/serializers"
	"github.com/thomas-maurice/cachou/storage"

	"github.com/stretchr/testify/assert"
)

type testInvalidStruct struct {
	ID    int64  `json:"id"`
	Field string `json:"field"`
}

type testValidStruct struct {
	ID    int64  `json:"id" cachou:"uid"`
	Field string `json:"field"`
}

func allTests(t *testing.T, cache *Cachou) {
	valid := testValidStruct{
		ID:    420,
		Field: "hello",
	}

	invalid := testInvalidStruct{
		ID:    69,
		Field: "hello",
	}

	cached, err := cache.Put(valid)
	assert.Nil(t, err)
	assert.True(t, cached)

	cached, err = cache.Put(invalid)
	assert.NotNil(t, err)
	assert.False(t, cached)
}

func TestCachou(t *testing.T) {
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)

	j := serializers.NewJSONSerializer()
	boltStorage, err := storage.NewBoltStorage(path.Join(tempDir, "bolt.db"))
	assert.Nil(t, err)

	memoryStorage := storage.NewMemoryStorage()

	caches := []*Cachou{
		NewCachou(j, boltStorage),
		NewCachou(j, memoryStorage),
	}

	for _, c := range caches {
		allTests(t, c)
	}

}
