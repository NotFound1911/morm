package morm

// Aggregate 聚合函数：AVG, MAX, MIN ...
type Aggregate struct {
	fn    string // 函数
	arg   string // 参数
	alias string // 别名
}

func (a Aggregate) selectable() {
}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
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
