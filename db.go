package morm

import (
	"database/sql"
	"github.com/NotFound1911/morm/internal/valuer"
	"github.com/NotFound1911/morm/model"
)

type DBOption func(*DB) error
type DB struct {
	r          model.Registry
	db         *sql.DB
	valCreator valuer.Creator
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
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}
func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:          model.NewRegistry(),
		db:         db,
		valCreator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		if err := opt(res); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// DBWithRegistry 使用自定义注册中心
func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) error {
		db.r = r
		return nil
	}
}

// DBUseReflectValuer 使用基于reflect的方法
func DBUseReflectValuer() DBOption {
	return func(db *DB) error {
		db.valCreator = valuer.NewReflectValue
		return nil
	}
}
