package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/NotFound1911/morm/model"
	"strings"
)

type builder struct {
	sqlBuilder strings.Builder
	args       []any
	model      *model.Model
	where      []Predicate

	table   TableReference // todo builder
	quoter  byte
	dialect Dialect
	core
}

func (b *builder) quote(name string) {
	b.sqlBuilder.WriteByte(b.quoter)
	b.sqlBuilder.WriteString(name)
	b.sqlBuilder.WriteByte(b.quoter)
}
func (b *builder) raw(r RawExpr) {
	b.sqlBuilder.WriteString(r.raw)
	if len(r.args) != 0 {
		b.addArgs(r.args...)
	}
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
		return b.buildColumn(exp.table, exp.name)
	case value: // 代表是列名，直接拼接列名
		b.sqlBuilder.WriteByte('?')
		b.args = append(b.args, exp.val)
	case RawExpr:
		b.raw(exp)
	case MathExpr:
		return b.buildBinaryExpr(binaryExpr(exp))
	case Subquery:
		return b.buildSubquery(exp, false)
	case SubqueryExpr:
		b.sqlBuilder.WriteString(exp.pred)
		b.sqlBuilder.WriteByte(' ')
		return b.buildSubquery(exp.s, false)
	case Predicate: // 代表查询条件
		return b.buildBinaryExpr(binaryExpr(exp))
	case Aggregate:
		return b.buildAggregate(exp, false)
	default:
		return errs.NewErrUnsupportedExpressionType(exp)
	}
	return nil
}
func (b *builder) buildBinaryExpr(e binaryExpr) error {
	err := b.buildSubExpr(e.left)
	if err != nil {
		return err
	}
	if e.opt != "" {
		b.sqlBuilder.WriteByte(' ')
		b.sqlBuilder.WriteString(e.opt.String())
	}
	if e.right != nil {
		b.sqlBuilder.WriteByte(' ')
		return b.buildSubExpr(e.right)
	}
	return nil
}
func (b *builder) buildSubExpr(subExpr Expression) error {
	switch sub := subExpr.(type) {
	case MathExpr:

	case binaryExpr:
	case Predicate:
		_ = b.sqlBuilder.WriteByte('(')
		if err := b.buildBinaryExpr(binaryExpr(sub)); err != nil {
			return err
		}
		_ = b.sqlBuilder.WriteByte(')')
	default:
		if err := b.buildExpression(sub); err != nil {
			return err
		}
	}
	return nil
}
func (b *builder) buildSubquery(table Subquery, useAlias bool) error {
	q, err := table.s.Build()
	if err != nil {
		return err
	}
	b.sqlBuilder.WriteByte('(')
	b.sqlBuilder.WriteString(q.SQL[:len(q.SQL)-1]) // 去掉;
	if len(q.Args) > 0 {
		b.addArgs(q.Args...)
	}
	b.sqlBuilder.WriteByte(')')
	if useAlias {
		b.sqlBuilder.WriteString(" AS ")
		b.quote(table.alias)
	}
	return nil
}

// buildColumn 构建列
// 如果 table 没有指定，用 model 来判断列是否存在
func (b *builder) buildColumn(table TableReference, fd string) error {
	var alias string
	if table != nil {
		alias = table.tableAlias()
	}
	if alias != "" {
		b.quote(alias)
		b.sqlBuilder.WriteByte('.')
	}
	colName, err := b.colName(table, fd)
	if err != nil {
		return err
	}
	b.quote(colName)

	return nil
}
func (b *builder) colName(table TableReference, fd string) (string, error) {
	switch tab := table.(type) {
	case nil:
		fdMeta, ok := b.model.FieldMap[fd]
		if !ok {
			return "", errs.NewErrUnknownField(fd)
		}
		return fdMeta.ColName, nil
	case Table:
		m, err := b.r.Get(tab.entity)
		if err != nil {
			return "", err
		}
		fdMeta, ok := m.FieldMap[fd]
		if !ok {
			return "", errs.NewErrUnknownField(fd)
		}
		return fdMeta.ColName, nil
	case Subquery:
		if len(tab.columns) > 0 {
			for _, col := range tab.columns {
				if col.selectedAlias() == fd {
					return fd, nil
				}
				if col.fieldName() == fd {
					return b.colName(col.target(), fd)
				}
			}
			return "", errs.NewErrUnknownField(fd)
		}
		return b.colName(tab.table, fd)
	default:
		return "", errs.NewErrUnsupportedExpressionType(tab)
	}
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
