# cachou: minimalist cache thingy

## Background
I wanted something like that to avoid un-necessary calls to external systems like Google Datastore and all these storage backends that can cost you money.
The idea is that you can pipe any Go struct into a cache, provided you have the right tags to make `get` calls quicker & cheaper.

## How it works
A `Cache` basically needs two things, a serializer and a storage engine. The `Serializer` is going to determine how your data is going to be cached, currently only JSON is supported but adding new ones is essentially trivial. The storage engine de termines *where* your data is going to be stored, right now you can use in memory storage, redis and bolt storage.

The point of this library is basically to be a middleware between the process comsuming data, and the source of truth itself. As an example you could have something like that

```go
package main

import (
	"log"

	"github.com/thomas-maurice/cachou"
	"github.com/thomas-maurice/cachou/serializers"
    "github.com/thomas-maurice/cachou/storage"
    "cloud.google.com/go/datastore"
)

type User struct {
	ID       int64  `json:"id" datastore:"id" cachou:"uid"` // this *must be unique*
	Username string `json:"username" datastore:"username"`
}

func saveUser(dsClient *datastore.Client, cache *cachou.Cachou, user User) error {
    k := datastore.IDKey("user", user.ID, nil)
	if _, err := s.dsClient.Put(context.Background(), k, user); err != nil {
		return err
	}

    
    _, err := cache.Put(user)
    return err
}

func getUser(dsClient *datastore.Client, cache *cachou.Cachou, id int64) (*User, error) {
    var user User
    found, err := cache.Get(&user, id)
    if err != nil {
        return nil, err
    }

    if found {
        return &user, nil
    }

    u := new(User)
	k := datastore.IDKey("user", id, nil)
	if err := dsClient.Get(context.Background(), k, u); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}

    _, err := cache.Put(user)
    if err != nil {
        return nil, err
    }
	return u, nil
}

func main() {
    dsClient, err := datastore.NewClient(ctx, "some-gcp-project")
	if err != nil {
		return nil, err
    }
    
	cache := cachou.NewCachou(serializers.NewJSONSerializer(), storage.NewMemoryStorage())

	u := User{
		ID:       69,
		Username: "thomas",
	}
}
```