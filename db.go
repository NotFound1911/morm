package morm

import (
	"context"
	"database/sql"
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/NotFound1911/morm/internal/valuer"
	"github.com/NotFound1911/morm/model"
)

type DBOption func(*DB) error
type DB struct {
	db *sql.DB
	core
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// Open 创建一个 DB 实例。
// 默认情况下，该 DB 将使用 MySQL 作为方言
// 如果你使用了其它数据库，可以使用 DBWithDialect 指定
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		core: core{
			dialect:    MySQL,
			r:          model.NewRegistry(),
			valCreator: valuer.NewUnsafeValue,
		},
		db: db,
	}
	for _, opt := range opts {
		if err := opt(res); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// DBWithDialect 使用自定义方言
func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) error {
		db.dialect = dialect
		return nil
	}
}

// DBWithRegistry 使用自定义注册中心
func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) error {
		db.r = r
		return nil
	}
}

// DBWithMiddleware 使用中间件
func DBWithMiddleware(ms ...Middleware) DBOption {
	return func(db *DB) error {
		db.ms = ms
		return nil
	}
}

// DBUseReflectValuer 使用基于reflect的方法
func DBUseReflectValuer() DBOption {
	return func(db *DB) error {
		db.valCreator = valuer.NewReflectValue
		return nil
	}
}

// BeginTx开启事务
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx, db: db}, nil
}

func (db *DB) DoTx(ctx context.Context,
	fn func(ctx context.Context, tx *Tx) error,
	opts *sql.TxOptions) (err error) { // err 是保留最新的错误
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		if e := recover(); e != nil || err != nil {
			if e != nil {
				err = errs.NewErrTxFuncFailed(e)
			}
			rE := tx.Rollback()
			if rE != nil {
				err = errs.NewErrTxRollbackFailed(rE)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				err = errs.NewErrTxCommitFailed(err)
			}
		}
	}()
	err = fn(ctx, tx)
	return err
}

func (db *DB) Close() error {
	return db.db.Close()
}
