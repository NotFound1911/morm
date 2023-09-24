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

// Predicate 是一个查询条件，开源通过组合的方式构建复杂的查询条件
type Predicate struct {
	left  Expression
	opt   opt
	right Expression
}

var _ Expression = &Predicate{}

func (Predicate) expr() {

}

func Not(p Predicate) Predicate {
	return Predicate{
		opt:   optNOT,
		right: p,
	}
}

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		opt:   optAND,
		right: r,
	}
}
func (p Predicate) Or(r Predicate) Predicate {
	return Predicate{
		left:  p,
		opt:   optOR,
		right: r,
	}
}
