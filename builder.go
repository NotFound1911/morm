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
		fd, ok := b.model.FieldMap[exp.name]
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		b.sqlBuilder.WriteByte('`')
		b.sqlBuilder.WriteString(fd.ColName)
		b.sqlBuilder.WriteByte('`')
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
	default:
		return fmt.Errorf("orm: not support the expression %v", exp)
	}
	return nil
}
