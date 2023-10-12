package morm

import "github.com/NotFound1911/morm/model"

type DBOption func(*DB) error
type DB struct {
	r model.Registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	db := &DB{
		r: model.NewRegistry(),
	}
	for _, opt := range opts {
		if err := opt(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}
