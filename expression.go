package morm

// RawExpr 原生表达式，不做任何处理
type RawExpr struct {
	raw  string
	args []interface{}
}

func (e RawExpr) fieldName() string {
	return ""
}

func (e RawExpr) target() TableReference {
	return nil
}

func (e RawExpr) selectedAlias() string {
	return ""
}

func (e RawExpr) expr() {

}

func (e RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: e,
	}
}

func Raw(expr string, args ...interface{}) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

type binaryExpr struct {
	left  Expression
	opt   opt
	right Expression
}

func (b binaryExpr) expr() {
}

type MathExpr binaryExpr

func (m MathExpr) expr() {
}

func (m MathExpr) Add(val interface{}) MathExpr {
	return MathExpr{
		left:  m,
		opt:   optADD,
		right: valueOf(val),
	}
}
func (m MathExpr) Multi(val interface{}) MathExpr {
	return MathExpr{
		left:  m,
		opt:   optMULTI,
		right: valueOf(val),
	}
}

// SubqueryExpr 这个谓词这种不是在所有的数据库里面都支持的
type SubqueryExpr struct {
	s    Subquery
	pred string
}

func (s SubqueryExpr) expr() {
}

func Any(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s:    sub,
		pred: "ANY",
	}
}
func All(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s:    sub,
		pred: "ALL",
	}
}

func Some(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s:    sub,
		pred: "SOME",
	}
}
