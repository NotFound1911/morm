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
	var sqlByte strings.Builder
	sqlByte.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T
		sqlByte.WriteByte('`')
		sqlByte.WriteString(reflect.TypeOf(t).Name())
		sqlByte.WriteByte('`')
	} else {
		sqlByte.WriteString(s.table)
	}
	sqlByte.WriteString(";")
	return &Query{
		SQL: sqlByte.String(),
	}, nil
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}
