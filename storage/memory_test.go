package storage

import "testing"

func TestMemoryStorage(t *testing.T) {
	kv := NewMemoryStorage()

	genericTest(t, kv)
}
