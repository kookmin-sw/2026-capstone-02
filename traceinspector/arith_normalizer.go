package traceinspector

import (
	"fmt"
	"traceinspector/imp"
)

// Given an integer equality/inequality, rewrite to canonical form
// 1. e1 <= e2 -> e1 - e2 <= 0  I named this as zero-rhs form
// 2. ax + by + c <= 0 where a and b are integer constants and x y are identifiers

// Given an integer (in)equality expression of the form `e1 ☉ e2“, convert to `e1 - (e2) ☉ 0`
func zero_rhs(expr imp.Expr) (imp.Expr, error) {
	switch expr_ty := expr.(type) {
	case *imp.EqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.EqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.NeqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.NeqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.LessthanExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.LessthanExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.GreaterthanExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.GreaterthanExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.LeqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.LeqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	case *imp.GeqExpr:
		sub_expr := imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: expr_ty.Rhs}
		return &imp.GeqExpr{Lhs: &sub_expr, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 0}}, nil
	default:
		return nil, fmt.Errorf("zero_rhs: Unsupported boolean expression %s", expr)
	}
}

// Represents a linear arithmetic polynomial ax ☉ by ☉ ... ☉ cz + C,
// variable_expr: ax ☉ by ☉ ... ☉ cz
// constant: C
type polynomial struct {
	variable_expr imp.Expr
	constant      int
}

// Build the normalized polynomial representation of the integer expression
func build_polynomial(expr imp.Expr) (polynomial, error) {
	accumulated := polynomial{}
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		accumulated.variable_expr = expr_ty
	case *imp.IntLitExpr:
		accumulated.constant += expr_ty.Value
	case *imp.ArrayLitExpr:
		accumulated.variable_expr = expr_ty
	case *imp.ArrayIndexExpr:
		accumulated.variable_expr = expr_ty
	case *imp.MakeArrayExpr:
		accumulated.variable_expr = expr_ty
	case *imp.LenExpr:
		accumulated.variable_expr = expr_ty
	case *imp.CallExpr:
		accumulated.variable_expr = expr_ty
	case *imp.AddExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return polynomial{}, err
		}
		if lhs_poly.variable_expr == nil {
			rhs_poly.constant += lhs_poly.constant
			return rhs_poly, nil
		} else if rhs_poly.variable_expr == nil {
			lhs_poly.constant += rhs_poly.constant
			return lhs_poly, nil
		} else {
			accumulated.variable_expr = &imp.AddExpr{Node: expr_ty.Node, Lhs: lhs_poly.variable_expr, Rhs: rhs_poly.variable_expr}
			accumulated.constant += lhs_poly.constant + rhs_poly.constant
		}
	case *imp.SubExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return polynomial{}, err
		}
		if lhs_poly.variable_expr == nil {
			rhs_poly.constant += lhs_poly.constant
			return rhs_poly, nil
		} else if rhs_poly.variable_expr == nil {
			lhs_poly.constant -= rhs_poly.constant
			return lhs_poly, nil
		} else {
			accumulated.variable_expr = &imp.SubExpr{Node: expr_ty.Node, Lhs: lhs_poly.variable_expr, Rhs: rhs_poly.variable_expr}
			accumulated.constant += lhs_poly.constant - rhs_poly.constant
		}
	case *imp.MulExpr:
		// For the case of multiplication
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return polynomial{}, err
		}
		if lhs_poly.variable_expr == nil && rhs_poly.variable_expr == nil {
			// both subexprs are constants
			accumulated.constant = lhs_poly.constant * rhs_poly.constant

			// Don't do "constant folding" for subexprs yet
			// } else if lhs_poly.variable_expr == nil {
			// 	// LHS is constant, but RHS isn't so LHS should be used as coefficient
			// 	accumulated.variable_expr = &imp.MulExpr{Node: expr_ty.Node, Lhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: lhs_poly.constant}, Rhs: expr_ty.Rhs}
			// } else if rhs_poly.variable_expr == nil {
			// 	// same goes for RHS
			// 	accumulated.variable_expr = &imp.MulExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: rhs_poly.constant}}
			// } else {
		} else {
			// if both are not constants, return the original
			accumulated.variable_expr = expr_ty
		}
	case *imp.DivExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return polynomial{}, err
		}
		if lhs_poly.variable_expr == nil && rhs_poly.variable_expr == nil {
			accumulated.constant = lhs_poly.constant / rhs_poly.constant
		} else {
			accumulated.variable_expr = expr_ty
		}
	case *imp.ModExpr:
		lhs_poly, err := build_polynomial(expr_ty.Lhs)
		if err != nil {
			return polynomial{}, err
		}
		rhs_poly, err := build_polynomial(expr_ty.Rhs)
		if err != nil {
			return polynomial{}, err
		}
		if lhs_poly.variable_expr == nil && rhs_poly.variable_expr == nil {
			accumulated.constant = lhs_poly.constant % rhs_poly.constant
		} else {
			accumulated.variable_expr = expr_ty
		}
	case *imp.NegExpr:
		sub_poly, err := build_polynomial(expr_ty.Subexpr)
		if err != nil {
			return polynomial{}, err
		}
		if sub_poly.variable_expr == nil {
			accumulated.constant -= sub_poly.constant
		} else {
			accumulated.variable_expr = &imp.NegExpr{Node: expr_ty.Node, Subexpr: sub_poly.variable_expr}
			accumulated.constant -= accumulated.constant
		}
	case *imp.ParenExpr:
		sub_poly, err := build_polynomial(expr_ty.Subexpr)
		if err != nil {
			return polynomial{}, err
		}
		if sub_poly.variable_expr == nil {
			accumulated.constant += sub_poly.constant
		} else {
			accumulated.variable_expr = &imp.ParenExpr{Node: expr_ty.Node, Subexpr: sub_poly.variable_expr}
			accumulated.constant += accumulated.constant
		}
	default:
		return polynomial{}, fmt.Errorf("build_polynomial: unsupported expressions %s", expr_ty)
	}
	return accumulated, nil
}

// Given an arbitrary integer expression, normalize to the form
// ax ☉ by ☉ ... ☉ cz + C, where C is an integer constant
func normalize_integer_expr(expr imp.Expr) (imp.Expr, error) {
	poly, err := build_polynomial(expr)
	if err != nil {
		return nil, err
	} else {
		return &imp.AddExpr{Node: imp.Node{Line_num: expr.GetLineNum()}, Lhs: poly.variable_expr, Rhs: &imp.IntLitExpr{Node: imp.Node{Line_num: expr.GetLineNum()}, Value: poly.constant}}, nil
	}
}
