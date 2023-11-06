package morm

type Deleter[T any] struct {
	builder

	core
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
func NewDeleter[T any](sess session) *Deleter[T] {
	c := sess.getCore()
	return &Deleter[T]{
		core: c,
		sess: sess,
		builder: builder{
			dialect: c.dialect,
			quoter:  c.dialect.quoter(),
		},
	}
}
