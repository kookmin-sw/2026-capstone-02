package domain

import (
	"traceinspector/imp"
)

type FilterQueryType int

const (
	FilterQueryType_Invalid FilterQueryType = iota
	FilterQueryType_Eq
	FilterQueryType_Neq
	FilterQueryType_Leq
	FilterQueryType_Geq
)

type FilterQuery struct {
	term_expr  imp.Expr
	query_type FilterQueryType
	rhs        imp.Expr
}
