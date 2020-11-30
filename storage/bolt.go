package storage

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

type BoltStorage struct {
	db *bolt.DB
}

func NewBoltStorage(fileName string) (*BoltStorage, error) {
	db, err := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &BoltStorage{
		db: db,
	}, nil
}

func (s *BoltStorage) Put(objType []byte, objUID []byte, objData []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// TODO (tmaurice): dont call that more than once
		b, err := tx.CreateBucketIfNotExists(objType)
		if err != nil {
			return err
		}

		return b.Put(objUID, objData)
	})
}

func (s *BoltStorage) Get(objType []byte, objUID []byte) ([]byte, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(objType)
		if b == nil {
			return fmt.Errorf("no such bucket: %s", string(objType))
		}

		data = b.Get(objUID)
		return nil
	})

	return data, err
}

func (s *BoltStorage) Del(objType []byte, objUID []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(objType)
		if b != nil {
			return fmt.Errorf("no such bucket: %s", string(objType))
		}

		return b.Delete(objUID)
	})

}
