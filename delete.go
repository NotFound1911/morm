package morm

import (
	"github.com/NotFound1911/morm/errors"
)

type Deleter[T any] struct {
	builder
	sess session
}

func (d *Deleter[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	d.model, err = d.r.Get(&t)
	if err != nil {
		return nil, err
	}
	d.sqlBuilder.WriteString("DELETE FROM ")
	if err = d.buildTable(d.table); err != nil {
		return nil, err
	}
	// 构造where
	if len(d.where) > 0 {
		d.sqlBuilder.WriteString(" WHERE ")
		if err := d.buildPredicates(d.where); err != nil {
			return nil, err
		}
	}
	d.sqlBuilder.WriteString(";")
	return &Query{
		SQL:  d.sqlBuilder.String(),
		Args: d.args,
	}, nil
}
func (d *Deleter[T]) buildTable(table TableReference) error {
	switch tab := table.(type) {
	case nil:
		d.quote(d.model.TableName)
	case Table:
		model, err := d.r.Get(tab.entity)
		if err != nil {
			return err
		}
		d.quote(model.TableName)
	default:
		return errs.NewErrUnsupportedExpressionType(tab)
	}
	return nil
}

// From accepts model definition
func (d *Deleter[T]) From(table TableReference) *Deleter[T] {
	d.table = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(ps ...Predicate) *Deleter[T] {
	d.where = ps
	return d
}
func NewDeleter[T any](sess session) *Deleter[T] {
	c := sess.getCore()
	return &Deleter[T]{
		sess: sess,
		builder: builder{
			core:    c,
			dialect: c.dialect,
			quoter:  c.dialect.quoter(),
		},
	}
}
