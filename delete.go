package morm

type Deleter[T any] struct {
	builder
}

func (d *Deleter[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	d.model, err = d.db.r.Get(&t)
	if err != nil {
		return nil, err
	}
	d.sqlBuilder.WriteString("DELETE FROM ")
	if d.table == "" {
		d.sqlBuilder.WriteByte('`')
		d.sqlBuilder.WriteString(d.model.TableName)
		d.sqlBuilder.WriteByte('`')
	} else {
		d.sqlBuilder.WriteString(d.table)
	}
	// 构造where
	if len(d.where) > 0 {
		d.sqlBuilder.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}
	}
	d.sqlBuilder.WriteString(";")
	return &Query{
		SQL:  d.sqlBuilder.String(),
		Args: d.args,
	}, nil
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(ps ...Predicate) *Deleter[T] {
	d.where = ps
	return d
}
func NewDeleter[T any](db *DB) *Deleter[T] {
	return &Deleter[T]{
		builder: builder{
			db: db,
		},
	}
}
