package bolt

import (
	"encoding/json"
	"log"

	bolt "go.etcd.io/bbolt"
)

var (
	db   *bolt.DB
	path = ".db"
)

func Open() {
	bolt, err := bolt.Open(path, 7777, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = bolt
}

func Close() {
	db.Close()
}

type Store struct {
	name string
}

func (s *Store) Create(id string, payload any) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))

		if err != nil {
			return err
		}

		err = b.Put([]byte(id), buf)
		return err
	})

	return err
}

func (s *Store) Exist(id string) bool {
	var payload []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		payload = b.Get([]byte(id))
		return nil
	})

	if err != nil {
		return false
	}

	return payload != nil
}

func (s *Store) FindOne(id string) any {
	var buf []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		buf = b.Get([]byte(id))
		return nil
	})

	if err != nil {
		return nil
	}

	var payload any
	err = json.Unmarshal(buf, &payload)
	if err != nil {
		return nil
	}

	return payload
}

func (s *Store) FindMany(callback func(string, []byte)) error {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			callback(string(k), v)
		}

		return nil
	})

	return err
}

func (s *Store) Update(id string, payload any) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		if err != nil {
			return err
		}

		err = b.Put([]byte(id), buf)
		return err
	})

	return err
}

func (s *Store) Delete(id string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.name))
		err := b.Delete([]byte(id))
		return err
	})

	return err
}

func New(name string) *Store {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	return &Store{name}
}

func Empty(name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		return err
	})

	return err
}
