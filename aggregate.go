package morm

// Aggregate 聚合函数：AVG, MAX, MIN ...
type Aggregate struct {
	fn    string // 函数
	arg   string // 参数
	alias string // 别名
}

func (a Aggregate) selectable() {
}
func (Aggregate) expr() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
	}
}
func (a Aggregate) EQ(arg any) Predicate { // = 等于
	return Predicate{
		left:  a,
		opt:   optEQ,
		right: exprOf(arg),
	}
}
func (a Aggregate) LT(arg any) Predicate { // < 小于
	return Predicate{
		left:  a,
		opt:   optLT,
		right: exprOf(arg),
	}
}
func (a Aggregate) GT(arg any) Predicate { // > 大于
	return Predicate{
		left:  a,
		opt:   optGT,
		right: exprOf(arg),
	}
}
func Avg(c string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: c,
	}
}

func Max(c string) Aggregate {
	return Aggregate{
		fn:  "MAX",
		arg: c,
	}
}

func Count(c string) Aggregate {
	return Aggregate{
		fn:  "COUNT",
		arg: c,
	}
}
func Sum(c string) Aggregate {
	return Aggregate{
		fn:  "SUM",
		arg: c,
	}
}
