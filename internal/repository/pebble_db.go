package repository

import (
	"github.com/cockroachdb/pebble"
	"encoding/json"
)

type PebbleDB struct {
	db *pebble.DB
}

func NewPebbleDB(path string) (*PebbleDB, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &PebbleDB{db: db}, nil
}

func (p *PebbleDB) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.db.Set([]byte(key), data, nil)
}

func (p *PebbleDB) Get(key string, value interface{}) error {
	data, closer, err := p.db.Get([]byte(key))
	if err != nil {
		return err
	}
	defer closer.Close()
	return json.Unmarshal(data, value)
}

func (p *PebbleDB) Delete(key string) error {
	return p.db.Delete([]byte(key), nil)
}