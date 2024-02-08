package morm

import (
	"context"
	"github.com/NotFound1911/morm/errors"
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

	sess session
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) From(table TableReference) *Selector[T] {
	s.table = table
	return s
}

// Build 构建query
func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	s.model, err = s.r.Get(&t)
	if err != nil {
		return nil, err
	}
	s.sqlBuilder.WriteString("SELECT ")
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sqlBuilder.WriteString(" FROM ")
	if err = s.buildTable(s.table); err != nil {
		return nil, err
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

func (s *Selector[T]) buildTable(table TableReference) error {
	switch tab := table.(type) {
	case nil:
		s.quote(s.model.TableName)
	case Table:
		model, err := s.r.Get(tab.entity)
		if err != nil {
			return err
		}
		s.quote(model.TableName)
		if tab.alias != "" {
			s.sqlBuilder.WriteString(" AS ")
			s.quote(tab.alias)
		}
	case Join:
		return s.buildJoin(tab)
	case Subquery:
		return s.buildSubquery(tab, true)
	default:
		return errs.NewErrUnsupportedExpressionType(tab)
	}
	return nil
}
func (s *Selector[T]) buildColumn(c Column, useAlias bool) error {
	if err := s.builder.buildColumn(c.table, c.name); err != nil {
		return err
	}
	if useAlias {
		s.buildAs(c.alias)
	}
	return nil
}
func (s *Selector[T]) buildJoin(table Join) error {
	s.sqlBuilder.WriteByte('(')
	if err := s.buildTable(table.left); err != nil {
		return err
	}
	s.sqlBuilder.WriteString(" ")
	s.sqlBuilder.WriteString(table.typ)
	s.sqlBuilder.WriteString(" ")
	if err := s.buildTable(table.right); err != nil {
		return err
	}
	if len(table.using) > 0 {
		s.sqlBuilder.WriteString(" USING (")
		for i, col := range table.using {
			if i > 0 {
				s.sqlBuilder.WriteByte(',')
			}
			if err := s.buildColumn(Column{name: col}, false); err != nil {
				return err
			}
		}
		s.sqlBuilder.WriteString(")")
	}
	if len(table.on) > 0 {
		s.sqlBuilder.WriteString(" ON ")
		if err := s.buildPredicates(table.on); err != nil {
			return err
		}
	}
	s.sqlBuilder.WriteByte(')')
	return nil
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func NewSelector[T any](sess session) *Selector[T] {
	c := sess.getCore()
	return &Selector[T]{
		sess: sess,
		builder: builder{
			core:    c,
			dialect: c.dialect,
			quoter:  c.dialect.quoter(),
		},
	}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	res := get[T](ctx, s.core, s.sess, &QueryContext{
		Builder: s,
		Type:    "SELECT",
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	res := getMultiHandler[T](ctx, s.core, s.sess, &QueryContext{
		Builder: s,
		Type:    "SELECT",
	})
	if res.Result != nil {
		return res.Result.([]*T), res.Err
	}
	return nil, res.Err
}

type Selectable interface {
	selectedAlias() string
	fieldName() string
	target() TableReference
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
		case RawExpr: //  表达式
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

func (s *Selector[T]) AsSubquery(alias string) Subquery {
	table := s.table
	if table == nil {
		table = TableOf(new(T))
	}
	return Subquery{
		s:       s,
		alias:   alias,
		table:   table,
		columns: s.columns,
	}
}
