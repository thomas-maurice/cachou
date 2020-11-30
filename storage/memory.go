package storage

import "fmt"

type MemoryStorage struct {
	objects map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		objects: make(map[string][]byte),
	}
}

func (s *MemoryStorage) key(objType, objUID []byte) string {
	return fmt.Sprint("%s#%s", string(objType), string(objUID))
}

func (s *MemoryStorage) Put(objType []byte, objUID []byte, objData []byte) error {
	s.objects[s.key(objType, objUID)] = objData
	return nil
}

func (s *MemoryStorage) Get(objType []byte, objUID []byte) ([]byte, error) {
	data, ok := s.objects[s.key(objType, objUID)]
	if !ok {
		return nil, nil
	}

	return data, nil
}

func (s *MemoryStorage) Del(objType []byte, objUID []byte) error {
	delete(s.objects, s.key(objType, objUID))
	return nil
}
