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

	// ErrInvalidTagContent 错误的标签设置
	ErrInvalidTagContent

	// ErrTooManyReturnedColumns 返回列过多
	ErrTooManyReturnedColumns

	// ErrInsertZeroRow 插入空值
	ErrInsertZeroRow

	// ErrUnsupportedAssignableType 不支持的 Assignable 表达式
	ErrUnsupportedAssignableType
)
