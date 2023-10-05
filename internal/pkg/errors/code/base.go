package code

// 通用: 基本错误
// Code must start with 1xxxxx
const (
	// ErrUnknown 未知错误
	ErrUnknown int = iota + 100001

	// ErrPointerOnly 只支持一级指针输入
	ErrPointerOnly

	// ErrUnknownField 返回代表未知字段的错误
	ErrUnknownField

	// ErrUnsupportedExpressionType 返回一个不支持该 expression 错误信息
	ErrUnsupportedExpressionType
)
