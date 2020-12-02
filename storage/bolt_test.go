package storage

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoltStorage(t *testing.T) {
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)

	boltStorage, err := NewBoltStorage(path.Join(tempDir, "bolt.db"))
	assert.Nil(t, err)

	defer boltStorage.Close()

	genericTest(t, boltStorage)
}
