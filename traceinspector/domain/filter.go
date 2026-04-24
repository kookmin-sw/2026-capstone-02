package domain

import (
	"traceinspector/algebra"
	"traceinspector/imp"
)

// Given a Simpleprop, arrange for target expr and return target's coeff, prop type, and the RHS
// The result is target ⊙ ±other + C, which is returned through (target_sign, prop_type, rhs), where
// target_sign denotes the sign(coeff) of target, prop type denotes ⊙, and rhs denotes the RHS imp.Expr ±x ± other
// If target does not exist, sign coeff is zero, prop type is invalid, and expr nil. Otherwise target_sign is always positive.
func arrange_for_expr(sp algebra.SimpleProp, target_expr_is_x bool) (algebra.SimplePropCoeff, FilterQueryType, imp.Expr) {
	var other_expr imp.Expr
	var target_coeff, other_coeff algebra.SimplePropCoeff
	switch target_expr_is_x {
	case true:
		// arrange for x
		target_coeff = sp.X_coeff
		other_expr = sp.Y_expr
		other_coeff = sp.Y_coeff
	case false:
		// arrange for y
		target_coeff = sp.Y_coeff
		other_expr = sp.X_expr
		other_coeff = sp.X_coeff
	}

	const_litnode := imp.IntLitExpr{Node: imp.Node{Line_num: 1337}, Value: sp.Constant}
	var rhs_expr imp.Expr
	var query_type FilterQueryType
	switch sp.Prop_type {
	case algebra.SimplePropType_Invalid:
		// invalid case
		return algebra.SimplePropCoeff_zero, FilterQueryType_Invalid, nil
	default:
		// ±target ±other ⊙ C
		// 1. Send other expr to the RHS, negating its coeff
		switch other_coeff.Negate() {
		case algebra.SimplePropCoeff_zero:
			// other doesn't exist, just emit coeff
			rhs_expr = &const_litnode
		case algebra.SimplePropCoeff_positive:
			rhs_expr = &imp.AddExpr{Node: imp.Node{Line_num: other_expr.GetLineNum()}, Lhs: &const_litnode, Rhs: other_expr}
		case algebra.SimplePropCoeff_negative:
			rhs_expr = &imp.SubExpr{Node: imp.Node{Line_num: other_expr.GetLineNum()}, Lhs: &const_litnode, Rhs: other_expr}
		}
	}
	switch sp.Prop_type {
	case algebra.SimplePropType_Invalid:
		query_type = FilterQueryType_Invalid
	case algebra.SimplePropType_Eq:
		query_type = FilterQueryType_Eq
	case algebra.SimplePropType_Neq:
		query_type = FilterQueryType_Neq
	case algebra.SimplePropType_Leq:
		query_type = FilterQueryType_Leq
	}

	// now we have the form ±target ⊙ C ∓other
	// if -target, negate RHS (and LEQ)
	switch target_coeff {
	case algebra.SimplePropCoeff_zero:
		// turns out target didn't exist!
		return algebra.SimplePropCoeff_zero, FilterQueryType_Invalid, nil
	case algebra.SimplePropCoeff_positive:
		// positive we do nothing
		return target_coeff, query_type, rhs_expr
	case algebra.SimplePropCoeff_negative:
		// -target ⊙ C ∓other
		// negate RHS and reason on ⊙
		rhs_expr = &imp.NegExpr{Node: imp.Node{Line_num: rhs_expr.GetLineNum()}, Subexpr: rhs_expr}
		if query_type == FilterQueryType_Leq {
			// -target <= C ∓other -> target >= -C ±other
			query_type = FilterQueryType_Geq
		}
		return algebra.SimplePropCoeff_positive, query_type, rhs_expr
	}
	return algebra.SimplePropCoeff_zero, FilterQueryType_Invalid, nil
}

// Given a SimpleProp, generate the filter queries for each subexpression assuming the prop is true.
// For example, x + y <= 3 will generate (<=, 3 - y) for x, and (<=, 3 - x) for y.
func Filter_true_query_simpleprop(sp algebra.SimpleProp) []FilterQuery {
	if sp.Y_coeff != algebra.SimplePropCoeff_zero {
		// the y term exists
		// arrange into form y ⊙ ± x + C
		//        term_expr--- | ----------Filterquery.rhs
		//                FilterQueryType
	}
	var ret []FilterQuery
	for_x_coeff, for_x_query, for_x_rhs := arrange_for_expr(sp, true)
	for_y_coeff, for_y_query, for_y_rhs := arrange_for_expr(sp, false)
	if for_x_coeff == algebra.SimplePropCoeff_positive {
		ret = append(ret, FilterQuery{Term_expr: sp.X_expr, Query_type: for_x_query, Rhs_expr: for_x_rhs})
	}
	if for_y_coeff == algebra.SimplePropCoeff_positive {
		ret = append(ret, FilterQuery{Term_expr: sp.Y_expr, Query_type: for_y_query, Rhs_expr: for_y_rhs})
	}
	return ret
}

// Like Filter_true_query, but for the negation of sp
func Filter_false_query_simpleprop(sp algebra.SimpleProp) []FilterQuery {
	return Filter_true_query_simpleprop(sp.Negate())
}
