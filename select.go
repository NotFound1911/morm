package morm

import (
	"context"
	"database/sql"
)

// Selector 构造select语句
type Selector[T any] struct {
	table string

	builder

	where []Predicate
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

// Build 构建query
func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	s.model, err = s.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
	s.sqlBuilder.WriteString("SELECT * FROM ")
	if s.table == "" {
		s.sqlBuilder.WriteByte('`')
		s.sqlBuilder.WriteString(s.model.TableName)
		s.sqlBuilder.WriteByte('`')
	} else {
		s.sqlBuilder.WriteString(s.table)
	}
	// 构造where
	if len(s.where) > 0 {
		s.sqlBuilder.WriteString(" WHERE ")
		if err := s.buildPredicates(s.where); err != nil {
			return nil, err
		}
	}
	s.sqlBuilder.WriteString(";")
	return &Query{
		SQL:  s.sqlBuilder.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		builder: builder{
			db: db,
		},
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := s.builder.db.db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	tmpl := new(T)
	meta, err := s.db.r.Get(tmpl)
	if err != nil {
		return nil, err
	}
	val := s.db.valCreator(tmpl, meta)
	err = val.SetColumns(rows)
	return tmpl, err
}

// todo
func (s *Selector[T]) GetMuti(ctx context.Context) ([]*T, error) {
	var db sql.DB
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		//todo
	}
	panic("implement me")
}
