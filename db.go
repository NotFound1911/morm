package morm

type DBOption func(*DB) error
type DB struct {
	r *registry
}

func NewDB(opts ...DBOption) (*DB, error) {
	db := &DB{
		r: &registry{},
	}
	for _, opt := range opts {
		if err := opt(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}
