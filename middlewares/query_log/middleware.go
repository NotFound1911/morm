package querylog

import (
	"context"
	"github.com/NotFound1911/morm"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(query string, args []any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("SQL: %s, args:%s \n", query, args)
		},
	}
}

// LogFunc 自定义日志函数
func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() morm.Middleware {
	return func(next morm.HanderFunc) morm.HanderFunc {
		return func(ctx context.Context, qc *morm.QueryContext) *morm.QueryResult {
			q, err := qc.Builder.Build() // 构造SQL和参数
			if err != nil {
				return &morm.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args)
			res := next(ctx, qc)
			return res
		}
	}
}
