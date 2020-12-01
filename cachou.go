package cachou

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	cachouTagName     = "cachou"
	cachouTagValueUID = "uid"
)

var (
	ErrObjectNotCacheable  = errors.New("object is not cacheable")                                              // Object is not cacheable
	ErrNoUniqueIdentifier  = errors.New("object does not possess a 'cachou:\"uid\"' tagged field")              // Missing proper tags
	ErrWrongIdentifierType = errors.New("wrong type of unique identifier, only ints and strings are supported") // Bad UID type
)

type Serializer interface {
	Serialize(interface{}) ([]byte, error)
	Deserialize([]byte, interface{}) error
}

type KVStore interface {
	// Arguments in order
	// * type name, converted to bytes
	// * uid of the given type
	// * serialized data
	Put([]byte, []byte, []byte) error
	// Arguments in order
	// * type name, converted to bytes
	// * uid of the given type
	// returns the serialized data
	Get([]byte, []byte) ([]byte, error)
	// Arguments in order
	// * type name, converted to bytes
	// * uid of the given type
	Del([]byte, []byte) error
}

type cacheableStruct struct {
	UIDField string
}

// Cachou is the main cache structure
type Cachou struct {
	structures map[string]*cacheableStruct
	serializer Serializer
	kvstore    KVStore
	objects    map[string][]byte
}

// NewCachou returns a new Cachou given the provided serializer and key-value storage backend
func NewCachou(s Serializer, kv KVStore) *Cachou {
	return &Cachou{
		structures: make(map[string]*cacheableStruct),
		serializer: s,
		kvstore:    kv,
	}
}

// scanObject is called the first time we have to deal with an unknown struct
// it will determine if the object is cacheable and which fields are interesting
// such as the uid. It will also ensure that the pk field is correct
func (c *Cachou) scanObject(object interface{}) error {
	objType := reflect.TypeOf(object)
	if objType.Kind() == reflect.Ptr {
		objType = reflect.TypeOf(object).Elem()
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		cacheTag, ok := field.Tag.Lookup(cachouTagName)
		if !ok {
			continue
		}
		if cacheTag == cachouTagValueUID {
			switch field.Type.Kind() {
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.String:
				c.structures[objType.String()] = &cacheableStruct{
					UIDField: field.Name,
				}

				return nil
			default:
				// not a supported type for a UID
				c.structures[objType.String()] = nil
				return ErrWrongIdentifierType
			}
		}
	}
	c.structures[objType.String()] = nil
	return ErrObjectNotCacheable
}

// getUID returns the []byte representation of the UID of the given struct
func (c *Cachou) getUID(structure *cacheableStruct, object interface{}) ([]byte, error) {
	objType := reflect.TypeOf(object)
	obj := reflect.ValueOf(object)
	if objType.Kind() == reflect.Ptr {
		objType = reflect.TypeOf(object).Elem()
		obj = obj.Elem()
	}

	return c.uidToBytes(obj.FieldByName(structure.UIDField))
}

// uidToBytes is a helper to convert a struct's UID to bytes for serialization
func (c *Cachou) uidToBytes(uid reflect.Value) ([]byte, error) {
	switch uid.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(fmt.Sprintf("%d", uid.Int())), nil
	case reflect.String:
		return []byte(uid.String()), nil
	}

	return nil, ErrWrongIdentifierType
}

func (c *Cachou) serialize(object interface{}) ([]byte, error) {
	return c.serializer.Serialize(object)
}

func (c *Cachou) deserialize(data []byte, object interface{}) error {
	return c.serializer.Deserialize(data, object)
}

// Put puts a value in the cache
func (c *Cachou) Put(object interface{}) (bool, error) {
	objType := reflect.TypeOf(object)
	if objType.Kind() == reflect.Ptr {
		objType = reflect.TypeOf(object).Elem()
	}

	if structure, ok := c.structures[objType.String()]; ok {
		if structure == nil {
			return false, ErrObjectNotCacheable
		}

		uid, err := c.getUID(structure, object)
		if err != nil {
			return false, err
		}

		b, err := c.serialize(object)
		if err != nil {
			return false, err
		}

		err = c.kvstore.Put([]byte(objType.String()), uid, b)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	// in case the struct is unknown to us yet
	err := c.scanObject(object)
	if err != nil {
		return false, err
	}

	return c.Put(object)
}

// Get retrieves an object from storage and passes it through the deserializer
func (c *Cachou) Get(object interface{}, uid interface{}) (bool, error) {
	objType := reflect.TypeOf(object)
	if objType.Kind() == reflect.Ptr {
		objType = reflect.TypeOf(object).Elem()
	}

	if structure, ok := c.structures[objType.String()]; ok {
		if structure == nil {
			// Not cacheable
			return false, ErrObjectNotCacheable
		} else {
			byteUID, err := c.uidToBytes(reflect.ValueOf(uid))
			if err != nil {
				return false, err
			}
			data, err := c.kvstore.Get([]byte(objType.String()), byteUID)
			if err != nil {
				return false, err
			}

			err = c.deserialize(data, object)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}

// Del deletes an element from the cache
func (c *Cachou) Del(object interface{}, uid interface{}) (bool, error) {
	objType := reflect.TypeOf(object)
	if objType.Kind() == reflect.Ptr {
		objType = reflect.TypeOf(object).Elem()
	}

	byteUID, err := c.uidToBytes(reflect.ValueOf(uid))
	if err != nil {
		return false, err
	}

	return true, c.kvstore.Del([]byte(objType.String()), byteUID)
}
