package morm

import (
	"github.com/NotFound1911/morm/errors"
)

var (
	MySQL   Dialect = &mysqlDialect{}
	SQLite3 Dialect = &sqlite3Dialect{}
)

type Dialect interface {
	quoter() byte
	buildUpsert(b *builder, odk *Upsert) error
}

type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildUpsert(b *builder, odk *Upsert) error {
	//TODO implement me
	panic("implement me")
}

type mysqlDialect struct {
	standardSQL
}

func (m *mysqlDialect) quoter() byte {
	return '`'
}
func (m *mysqlDialect) buildUpsert(b *builder, odk *Upsert) error {
	b.sqlBuilder.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i, a := range odk.assigns {
		if i > 0 {
			b.sqlBuilder.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			fd, ok := b.model.FieldMap[assign.name]
			if !ok {
				return errs.NewErrUnknownField(assign.name)
			}
			b.quote(fd.ColName)
			b.sqlBuilder.WriteString("=VALUES(")
			b.quote(fd.ColName)
			b.sqlBuilder.WriteByte(')')
		case Assignment:
			err := b.buildColumn(nil, assign.name)
			if err != nil {
				return err
			}
			b.sqlBuilder.WriteString("=?")
			b.addArgs(assign.val)
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}

type sqlite3Dialect struct {
	standardSQL
}

func (s *sqlite3Dialect) quoter() byte {
	return '`'
}
func (s *sqlite3Dialect) buildUpsert(b *builder, odk *Upsert) error {
	b.sqlBuilder.WriteString(" ON CONFLICT")
	if len(odk.conflictColumns) > 0 {
		b.sqlBuilder.WriteByte('(')
		for i, col := range odk.conflictColumns {
			if i > 0 {
				b.sqlBuilder.WriteByte(',')
			}
			if err := b.buildColumn(nil, col.name); err != nil {
				return err
			}
		}
		b.sqlBuilder.WriteByte(')')
	}
	b.sqlBuilder.WriteString(" DO UPDATE SET ")
	for i, a := range odk.assigns {
		if i > 0 {
			b.sqlBuilder.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			fd, ok := b.model.FieldMap[assign.name]
			if !ok {
				return errs.NewErrUnknownField(assign.name)
			}
			b.quote(fd.ColName)
			b.sqlBuilder.WriteString("=excluded.")
			b.quote(fd.ColName)

		case Assignment:
			err := b.buildColumn(assign.table, assign.name)
			if err != nil {
				return err
			}
			b.sqlBuilder.WriteString("=?")
			b.addArgs(assign.val)
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}
