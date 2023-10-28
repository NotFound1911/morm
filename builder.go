package morm

import (
	"fmt"
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	model "github.com/NotFound1911/morm/model"
	"strings"
)

type builder struct {
	sqlBuilder strings.Builder
	args       []any
	model      *model.Model
	db         *DB
	where      []Predicate
	columns    []Selectable
	table      string
}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}
	if err := b.buildExpression(p); err != nil {
		return err
	}
	return nil
}
func (b *builder) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column: // 代表是列名，直接拼接列名
		return b.buildColumn(exp, false)
	case value: // 代表是列名，直接拼接列名
		b.sqlBuilder.WriteByte('?')
		b.args = append(b.args, exp.val)
	case Predicate: // 代表查询条件
		_, lp := exp.left.(Predicate) // 判断是否是查询条件
		if lp {
			b.sqlBuilder.WriteByte('(')
		}
		// 递归
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			b.sqlBuilder.WriteByte(')')
		}
		b.sqlBuilder.WriteByte(' ')
		b.sqlBuilder.WriteString(exp.opt.String())
		b.sqlBuilder.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			b.sqlBuilder.WriteByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			b.sqlBuilder.WriteByte(')')
		}
	case Aggregate:
		return b.buildAggregate(exp, false)
	default:
		return fmt.Errorf("orm: not support the expression %v", exp)
	}
	return nil
}

// 构建筛选列
func (b *builder) buildColumns() error {
	if len(b.columns) == 0 {
		b.sqlBuilder.WriteByte('*')
		return nil
	}
	for i, col := range b.columns {
		if i > 0 {
			b.sqlBuilder.WriteByte(',')
		}
		switch val := col.(type) {
		case Column: // 列
			if err := b.buildColumn(val, true); err != nil {
				return err
			}
		case Aggregate: // 聚合
			if err := b.buildAggregate(val, true); err != nil {
				return err
			}
		case Expr: //  表达式
			b.sqlBuilder.WriteString(val.raw)
			if len(val.args) != 0 {
				b.addArgs(val.args...)
			}
		}
	}
	return nil
}

// 构建列
func (b *builder) buildColumn(val Column, useAlias bool) error {
	b.sqlBuilder.WriteByte('`')
	fd, ok := b.model.FieldMap[val.name]
	if !ok {
		return errs.NewErrUnknownField(val.name)
	}
	b.sqlBuilder.WriteString(fd.ColName)
	b.sqlBuilder.WriteByte('`')
	if useAlias {
		b.buildAs(val.alias)
	}
	return nil
}

// 构建聚合
func (b *builder) buildAggregate(val Aggregate, useAlias bool) error {
	b.sqlBuilder.WriteString(val.fn)
	b.sqlBuilder.WriteString("(`")
	fd, ok := b.model.FieldMap[val.arg]
	if !ok {
		return errs.NewErrUnknownField(val.arg)
	}
	b.sqlBuilder.WriteString(fd.ColName)
	b.sqlBuilder.WriteString("`)")
	if useAlias {
		b.buildAs(val.alias)
	}
	return nil
}

// 构建别名
func (b *builder) buildAs(alias string) {
	if alias != "" {
		b.sqlBuilder.WriteString(" AS ")
		b.sqlBuilder.WriteByte('`')
		b.sqlBuilder.WriteString(alias)
		b.sqlBuilder.WriteByte('`')
	}
}

func (b *builder) addArgs(args ...any) {
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}
	b.args = append(b.args, args...)
}
