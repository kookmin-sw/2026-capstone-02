package traceinspector

import (
	"fmt"
	"traceinspector/imp"
)

// Represent an inequality of the form ±x ±y <= C and ±x <= C, where x and y are variables.
// Also supports the form .
type SimpleInequality struct {
	x_varname string // variable name
	x_is_pos  bool   // whether coefficient of x is positive
	y_exists  bool   // whether the y variable exists
	y_varname string // variable name
	y_is_pos  bool   // whether coefficient of y is positive
	constant  int    // constant value
}

func (ieq SimpleInequality) String() string {
	var x_sign string
	var y_sign string = "+"
	if !ieq.x_is_pos {
		x_sign = "-"
	}
	if !ieq.y_is_pos {
		y_sign = "-"
	}
	if ieq.y_exists {
		return fmt.Sprintf("%s%s %s %s <= %d", x_sign, ieq.x_varname, y_sign, ieq.y_varname, ieq.constant)
	} else {
		return fmt.Sprintf("%s%s <= %d", x_sign, ieq.x_varname, ieq.constant)
	}
}

// Given an imp.Expr, check if the expr is of the form +-var. Only used within imp_expr_to_simp_inequality
func _check_if_var(expr imp.Expr) (var_name string, is_negative bool) {
	var_name = ""
	is_negative = false
	switch expr_ty := expr.(type) {
	case *imp.NegExpr:
		is_negative = true
		var_name_res, is_negative_res := _check_if_var(expr_ty.Subexpr)
		is_negative = is_negative != is_negative_res // xor
		var_name = var_name_res
	case *imp.VarExpr:
		var_name = expr_ty.Name
	}
	return
}

// Also verify that an expression is either a variable or a negation of it
func _check_binary_expr(expr imp.Expr) (string, bool, string, bool) {
	switch expr_ty := expr.(type) {
	case *imp.AddExpr:
		x_varname, x_is_neg := _check_if_var(expr_ty.Lhs)
		y_varname, y_is_neg := _check_if_var(expr_ty.Rhs)
		return x_varname, x_is_neg, y_varname, y_is_neg
	case *imp.ParenExpr:
		return _check_binary_expr(expr_ty.Subexpr)
	}
	return "", false, "", false
}

// Given an imp leq expression, try and convert the expression into a SimpleInequality.
// Returns SimpleInequality, and a boolean indicating whether the conversion was possible.
// Very naive and lazy implementation btw
func imp_expr_to_simp_inequality(expr imp.Expr) (SimpleInequality, bool) {
	switch expr_ty := expr.(type) {
	case *imp.LessthanExpr:
		// convert to leq
		// lhs < rhs -> lhs <= rhs - 1
		return imp_expr_to_simp_inequality(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Lhs, Rhs: &imp.SubExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: &imp.IntLitExpr{Node: expr_ty.Node, Value: 1}}})
	case *imp.GreaterthanExpr:
		// lhs > rhs -> rhs < lhs
		return imp_expr_to_simp_inequality(&imp.LessthanExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.GeqExpr:
		// lhs >= rhs -> rhs <= lhs
		return imp_expr_to_simp_inequality(&imp.LeqExpr{Node: expr_ty.Node, Lhs: expr_ty.Rhs, Rhs: expr_ty.Lhs})
	case *imp.LeqExpr:
		// move all terms to lhs
		zero_expr, err := zero_rhs(expr)
		if err != nil {
			return SimpleInequality{}, false
		}
		zero_expr_leq, is_leq_expr := zero_expr.(*imp.LeqExpr)
		if !is_leq_expr {
			return SimpleInequality{}, false
		}

		// pull constants out of LHS by representing LHS as Polynomial struct
		lhs_poly, err := build_polynomial(convert_subtraction_to_neg(zero_expr_leq.Lhs, false))
		// fmt.Println(expr, "->", zero_expr_leq.Lhs, "->", convert_subtraction_to_neg(zero_expr_leq.Lhs, false), "||", lhs_poly.variable_expr, lhs_poly.constant)
		created_ineq := SimpleInequality{}
		created_ineq.constant = -lhs_poly.constant // send constant to other side of leq

		// check if the polynomial is the form `±x + C`
		single_varname, single_var_is_neg := _check_if_var(lhs_poly.variable_expr)
		if single_varname != "" {
			created_ineq.x_varname = single_varname
			created_ineq.y_exists = false
			created_ineq.x_is_pos = !single_var_is_neg
			return created_ineq, true
		}

		// check if the polynomial is the form `±x ±y`
		x_varname, x_is_neg, y_varname, y_is_neg := _check_binary_expr(lhs_poly.variable_expr)

		if x_varname != "" && y_varname != "" {
			created_ineq.x_varname = x_varname
			created_ineq.x_is_pos = !x_is_neg
			created_ineq.y_exists = true
			created_ineq.y_varname = y_varname
			created_ineq.y_is_pos = !y_is_neg
			return created_ineq, true
		}
	}
	return SimpleInequality{}, false
}
