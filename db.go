package Loomup

import (
	"appengine"
	"appengine/channel"
	"appengine/datastore"
	"appengine/memcache"
)

// Building are stored in the datastore to be the parent entity of many Clients,
// keeping all the participants in a particular chat in the same entity group.

// Building represents a Pysical buidling or place sourrended by a geo lat/long.
type Building struct {
	Name string
	Address string
	Occupants int
}
type Vertex struct {
	lat, lng float64
}
func (r *Building) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Building", r.Name, 0, nil)
	
}
// Person is an Occupant in a Building.
type Client struct {
	ClientID string // the channel Client ID
}

// AddClient puts a Client record to the datastore with the Building as its
// parent, creates a channel and returns the channel token.
func (r *Building) AddClient(c appengine.Context, id string) (string, error) {
	key := datastore.NewKey(c, "Client", id, 0, r.Key(c))
	client := &Client{id}
	_, err := datastore.Put(c, key, client)
	if err != nil {
		return "", err
	}

	// Purge the now-invalid cache record (if it exists).
	memcache.Delete(c, r.Name)
	return channel.Create(c, id)
}

func (r *Building) Send(c appengine.Context, message string) error {
	var clients []Client

	_, err := memcache.JSON.Get(c, r.Name, &clients)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	if err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("Client").Ancestor(r.Key(c))
		_, err = q.GetAll(c, &clients)
		if err != nil {
			return err
		}
		err = memcache.JSON.Set(c, &memcache.Item{
			Key: r.Name, Object: clients,
		})
		if err != nil {
			return err
		}
	}

	for _, client := range clients {
		err = channel.Send(c, client.ClientID, message)
		if err != nil {
			c.Errorf("sending %q: %v", message, err)
		}
	}

	return nil
}

// getBuilding fetches a Building by name from the datastore,
// creating it if it doesn't exist already.
func getBuilding(c appengine.Context, name string) (*Building, error) {
	building := &Building{name, "", 0}

	fn := func(c appengine.Context) error {
		err := datastore.Get(c, building.Key(c), building)
		if err == datastore.ErrNoSuchEntity {
			_, err = datastore.Put(c, building.Key(c), building)
		}
		return err
	}

	return building, datastore.RunInTransaction(c, fn, nil)
}
