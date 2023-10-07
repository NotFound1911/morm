package errs

import (
	"fmt"
	"github.com/NotFound1911/morm/internal/pkg/errors/code"
)

type withCode struct {
	err  error
	code int
}

func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:  fmt.Errorf(format, args...),
		code: code,
	}
}

func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

func NewErrUnknown(exp any) error {
	return WithCode(code.ErrUnknown, fmt.Sprintf("morm 未知错误:%+v", exp))
}

func NewErrPointerOnly(exp any) error {
	return WithCode(code.ErrPointerOnly, fmt.Sprintf("morm 只支持一级指针作为输入, 例如 *User, %+v 不支持 ", exp))
}

func NewErrUnknownField(exp any) error {
	return WithCode(code.ErrUnknownField, fmt.Sprintf("morm 未知字段:%+v", exp))
}

func NewErrUnsupportedExpressionType(exp any) error {
	return WithCode(code.ErrUnsupportedExpressionType, fmt.Sprintf("morm 不支持表达式:%+v", exp))
}

func NewErrInvalidTagContent(exp any) error {
	return WithCode(code.ErrInvalidTagContent, fmt.Sprintf("morm 错误的标签设置:%+v", exp))
}
