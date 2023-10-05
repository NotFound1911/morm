package morm

import (
	"fmt"
	"reflect"
	"strings"
)

// Selector 构造select语句
type Selector[T any] struct {
	table      string
	sqlBuilder strings.Builder
	args       []any
	where      []Predicate
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

// Build 构建query
func (s *Selector[T]) Build() (*Query, error) {
	s.sqlBuilder.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T
		s.sqlBuilder.WriteByte('`')
		s.sqlBuilder.WriteString(reflect.TypeOf(t).Name())
		s.sqlBuilder.WriteByte('`')
	} else {
		s.sqlBuilder.WriteString(s.table)
	}
	// 构造where
	if len(s.where) > 0 {
		s.sqlBuilder.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}
	s.sqlBuilder.WriteString(";")
	return &Query{
		SQL:  s.sqlBuilder.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column: // 代表是列名，直接拼接列名
		s.sqlBuilder.WriteByte('`')
		s.sqlBuilder.WriteString(exp.name)
		s.sqlBuilder.WriteByte('`')
	case value: // 代表是列名，直接拼接列名
		s.sqlBuilder.WriteByte('?')
		s.args = append(s.args, exp.val)
	case Predicate: // 代表查询条件
		_, lp := exp.left.(Predicate) // 判断是否是查询条件
		if lp {
			s.sqlBuilder.WriteByte('(')
		}
		// 递归
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sqlBuilder.WriteByte(')')
		}
		s.sqlBuilder.WriteByte(' ')
		s.sqlBuilder.WriteString(exp.opt.String())
		s.sqlBuilder.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			s.sqlBuilder.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sqlBuilder.WriteByte(')')
		}
	default:
		return fmt.Errorf("orm: not support the expression %v", exp)
	}
	return nil
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}
