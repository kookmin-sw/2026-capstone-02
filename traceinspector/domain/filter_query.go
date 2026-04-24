package domain

import (
	"fmt"
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

func (qty FilterQueryType) String() string {
	switch qty {
	case FilterQueryType_Eq:
		return "=_f"
	case FilterQueryType_Neq:
		return "!=_f"
	case FilterQueryType_Leq:
		return "<=_f"
	case FilterQueryType_Geq:
		return ">=_f"
	}
	return "INVALID_QUERYTYPE"
}

type FilterQuery struct {
	Term_expr  imp.Expr
	Query_type FilterQueryType
	Rhs_expr   imp.Expr
}

func (query FilterQuery) String() string {
	return fmt.Sprintf("Filter(%s %s %s)", query.Term_expr, query.Query_type, query.Rhs_expr)
}
