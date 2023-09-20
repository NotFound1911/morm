package morm

import (
	"reflect"
	"strings"
)

// Selector 构造select语句
type Selector[T any] struct {
	table string
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

// Build 构建query
func (s *Selector[T]) Build() (*Query, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T
		sqlBuilder.WriteByte('`')
		sqlBuilder.WriteString(reflect.TypeOf(t).Name())
		sqlBuilder.WriteByte('`')
	} else {
		sqlBuilder.WriteString(s.table)
	}
	sqlBuilder.WriteString(";")
	return &Query{
		SQL: sqlBuilder.String(),
	}, nil
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}
