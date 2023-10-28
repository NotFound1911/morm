package morm

// 操作符
type opt string

const (
	optEQ  = "="
	optLT  = "<"
	optGT  = ">"
	optAND = "AND"
	optOR  = "OR"
	optNOT = "NOT"
)

func (o opt) String() string {
	return string(o)
}

type Expression interface {
	expr()
}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return valueOf(exp)
	}
}

// Predicate 是一个查询条件，通过组合的方式构建复杂的查询条件
type Predicate struct {
	left  Expression // 左边查询条件
	opt   opt        // 操作符
	right Expression // 右边查询条件
}

var _ Expression = &Predicate{}

func (Predicate) expr() {

}

// Not Not(查询条件 p)
func Not(p Predicate) Predicate {
	return Predicate{
		opt:   optNOT,
		right: p,
	}
}

// And (查询条件 p)And(查询条件 r)
func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		opt:   optAND,
		right: r,
	}
}

// Or (查询条件 p)Or(查询条件 r)
func (p Predicate) Or(r Predicate) Predicate {
	return Predicate{
		left:  p,
		opt:   optOR,
		right: r,
	}
}
