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

	// valid put
	cached, err := cache.Put(valid)
	assert.Nil(t, err)
	assert.True(t, cached)

	// invalid put
	cached, err = cache.Put(invalid)
	assert.NotNil(t, err)
	assert.False(t, cached)

	// valid get
	var cachedStruct testValidStruct
	gotten, err := cache.Get(&cachedStruct, 420)
	assert.Nil(t, err)
	assert.True(t, gotten)

	// invalid get
	var invalidCachedStruct testInvalidStruct
	gotten, err = cache.Get(&invalidCachedStruct, 420)
	assert.NotNil(t, err)
	assert.False(t, gotten)

	// valid delete
	deleted, err := cache.Del(testValidStruct{}, 420)
	assert.Nil(t, err)
	assert.True(t, deleted)
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
