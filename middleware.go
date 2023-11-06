package morm

import (
	"context"
	"github.com/NotFound1911/morm/model"
)

type QueryContext struct {
	// Type 操作类型 SELECT UPDATE DELETE INSERT
	Type string

	Builder QueryBuilder
	Model   *model.Model
}

type QueryResult struct {
	// Result 在不同的处理下可能是不同的 可能是单个结果，可能是多个结果
	Result any
	Err    error
}

type Middleware func(next HanderFunc) HanderFunc

type HanderFunc func(ctx context.Context, qc *QueryContext) *QueryResult
