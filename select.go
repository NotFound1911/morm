package morm

import (
	"context"
	"database/sql"
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
)

// Selector 构造select语句
type Selector[T any] struct {
	builder
	orderBys []OrderBy
	limit    int
	offset   int
	groupBys []Column
	having   []Predicate
	columns  []Selectable
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
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
	s.sqlBuilder.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sqlBuilder.WriteString(" FROM ")
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
	// 构造order by
	if len(s.orderBys) > 0 {
		s.sqlBuilder.WriteString(" ORDER BY ")
		for i, order := range s.orderBys {
			if i > 0 {
				s.sqlBuilder.WriteByte(',')
			}
			fd, ok := s.model.FieldMap[order.col]
			if !ok {
				return nil, errs.NewErrUnknownField(order.col)
			}
			s.sqlBuilder.WriteByte('`')
			s.sqlBuilder.WriteString(fd.ColName)
			s.sqlBuilder.WriteByte('`')
			s.sqlBuilder.WriteByte(' ')
			s.sqlBuilder.WriteString(order.fun)
		}
	}
	if s.limit > 0 {
		s.sqlBuilder.WriteString(" LIMIT ?")
		s.args = append(s.args, s.limit)
	}
	if s.offset > 0 {
		s.sqlBuilder.WriteString(" OFFSET ?")
		s.args = append(s.args, s.offset)
	}
	// group by
	if err = s.buildGroupBy(); err != nil {
		return nil, err
	}
	// having
	if err = s.buildHaving(); err != nil {
		return nil, err
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
			db:      db,
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
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

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := s.builder.db.db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}
	tmpls := make([]*T, 0, 0)
	tmpl := new(T)
	meta, err := s.db.r.Get(tmpl)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tmpl := new(T)
		val := s.db.valCreator(tmpl, meta)
		if err := val.SetColumns(rows); err != nil {
			return nil, err
		}
		tmpls = append(tmpls, tmpl)
	}
	if len(tmpls) == 0 {
		return nil, sql.ErrNoRows
	}
	return tmpls, nil
}

type Selectable interface {
	selectable()
}

// GroupBy 设置 group by 子句
func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T] {
	s.groupBys = cols
	return s
}
func (s *Selector[T]) buildGroupBy() error {
	if len(s.groupBys) > 0 {
		s.sqlBuilder.WriteString(" GROUP BY ")
		for i, group := range s.groupBys {
			if i > 0 {
				s.sqlBuilder.WriteByte(',')
			}
			if err := s.buildColumn(group, false); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	s.having = ps
	return s
}
func (s *Selector[T]) buildHaving() error {
	if len(s.having) > 0 {
		s.sqlBuilder.WriteString(" HAVING ")
		if err := s.buildPredicates(s.having); err != nil {
			return err
		}
	}
	return nil
}
func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}

func (s *Selector[T]) OrderBy(orderBys ...OrderBy) *Selector[T] {
	s.orderBys = orderBys
	return s
}

// 构建筛选列
func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sqlBuilder.WriteByte('*')
		return nil
	}
	for i, col := range s.columns {
		if i > 0 {
			s.sqlBuilder.WriteByte(',')
		}
		switch val := col.(type) {
		case Column: // 列
			if err := s.buildColumn(val, true); err != nil {
				return err
			}
		case Aggregate: // 聚合
			if err := s.buildAggregate(val, true); err != nil {
				return err
			}
		case Expr: //  表达式
			s.sqlBuilder.WriteString(val.raw)
			if len(val.args) != 0 {
				s.addArgs(val.args...)
			}
		}
	}
	return nil
}

type OrderBy struct {
	col string
	fun string
}

func Asc(col string) OrderBy { // 顺序
	return OrderBy{
		col: col,
		fun: "ASC",
	}
}

func Desc(col string) OrderBy { // 逆序
	return OrderBy{
		col: col,
		fun: "DESC",
	}
}
