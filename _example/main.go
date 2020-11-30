package main

import (
	"log"

	"github.com/thomas-maurice/cachou"
	"github.com/thomas-maurice/cachou/serializers"
	"github.com/thomas-maurice/cachou/storage"
)

type User struct {
	ID       int64  `json:"id" cachou:"uid"`
	Username string `json:"username"`
}

func main() {
	cache := cachou.NewCachou(serializers.NewJSONSerializer(), storage.NewMemoryStorage())

	u := User{
		ID:       69,
		Username: "thomas",
	}

	cached, err := cache.Put(u)
	if err != nil {
		panic(err)
	}

	log.Println("cached: ", cached)
	// now lets retrieve it from the cache

	var u2 User
	gotten, err := cache.Get(&u2, 69)
	log.Println("cache hit: ", gotten)
	log.Println(u2)

	// Now with an actual physical db
	db, err := storage.NewBoltStorage("foo.bolt-db")
	if err != nil {
		panic(err)
	}
	cache = cachou.NewCachou(serializers.NewJSONSerializer(), db)

	u = User{
		ID:       69,
		Username: "thomas",
	}

	cached, err = cache.Put(u)
	if err != nil {
		panic(err)
	}

	log.Println("cached: ", cached)
	// now lets retrieve it from the cache

	gotten, err = cache.Get(&u2, 69)
	log.Println("cache hit: ", gotten)
	log.Println(u2)

}
