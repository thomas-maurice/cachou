package storage

import (
	"testing"

	"github.com/thomas-maurice/cachou"

	"github.com/stretchr/testify/assert"
)

func genericTest(t *testing.T, storage cachou.KVStore) {
	k := []byte("type")
	uid := []byte("foo")
	value := []byte("value")

	err := storage.Put(k, uid, value)
	assert.Nil(t, err)

	v, err := storage.Get(k, uid)
	assert.Nil(t, err)
	assert.Equal(t, v, value)

	v, err = storage.Get(k, []byte("6969"))
	assert.Nil(t, err)
	assert.Equal(t, []uint8([]byte(nil)), v)

	err = storage.Del(k, uid)
	assert.Nil(t, err)

	v, err = storage.Get(k, uid)
	assert.Nil(t, err)
	assert.Equal(t, []uint8([]byte(nil)), v)
}
