package traceinspector

import (
	"fmt"
	"strings"
	"traceinspector/algebra"
	"traceinspector/domain"
	"traceinspector/imp"
)

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_location: node location to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
type AbstractState[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] struct {
	node_location CFGNodeLocation
	abstract_mem  AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]
}

func (state AbstractState[IntDomainImpl, ArrayDomainImpl]) String() string {
	var ret []string
	for key, val := range state.abstract_mem {
		ret = append(ret, fmt.Sprintf("%s : %s", key, val))
	}
	return fmt.Sprintf("%s - {%s}", state.node_location, strings.Join(ret, ", "))
}

func (state AbstractState[IntDomainImpl, ArrayDomainImpl]) Clone() AbstractState[IntDomainImpl, ArrayDomainImpl] {
	new_st := AbstractState[IntDomainImpl, ArrayDomainImpl]{}
	new_st.node_location = state.node_location
	new_st.abstract_mem = state.abstract_mem.Clone()
	return new_st
}

// Step: Given an input state (l, m^#), execute the abstract step relation for l under memory state m^#, and
// Return the subsequent states {(l', m^#')} ∈ P(L * M^#)
type AbstractSemantics[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] interface {
	Step(AbstractState[IntDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, ArrayDomainImpl]
}

// Abstract transition semantics for Imp wrt to arbitrary abstract domain impelmentations

// ImpFunctionInterpreter performs abstract interpretation of a function body from a given initial state. The
// interpreter performs interpretation until it collects the fixpoint semantics for the function body, and hence the
// return value. The interpreter will spawn another ImpFunctionInterpreter in the case a function call is invoked.
type ImpFunctionInterpreter[IntDomainImpl domain.IntegerDomain[IntDomainImpl], ArrayDomainImpl ArrayDomain[IntDomainImpl, ArrayDomainImpl]] struct {
	func_cfg_map        FunctionCFGMap
	func_name           imp.ImpFunctionName
	func_info_map       imp.ImpFunctionMap
	abstract_mem        *FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl] // joined global state
	intdomain_default   IntDomainImpl                                        // an instantiation of the integer domain impl
	booldomain_default  domain.BoolDomain                                    // an instantiation of the boolean domain
	arraydomain_default ArrayDomainImpl                                      // an instantiation of the array domain impl
}

// Compute the abstract value of an expression expr under an abstract memory state abs_mem
func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) Eval_expr(node_location CFGNodeLocation, expr imp.Expr, abs_mem AbstractNodeMem[IntDomainImpl, ArrayDomainImpl]) AbstractValue[IntDomainImpl, ArrayDomainImpl] {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		return abs_mem[expr_ty.Name]
	case *imp.IntLitExpr:
		intdom_result := interpreter.intdomain_default.From_IntLitExpr(*expr_ty)
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: intdom_result}
	case *imp.BoolLitExpr:
		booldom_result := interpreter.booldomain_default.From_BoolLitExpr(*expr_ty)
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: booldom_result}
	case *imp.StringLitExpr:
		// TODO do something about strings
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
	case *imp.NegExpr:
		subexpr_val := interpreter.Eval_expr(node_location, expr_ty.Subexpr, abs_mem)
		if subexpr_val.domain_kind != IntDomainKind {
			write_error(node_location, fmt.Sprintf("Result of arithmetic negation returned %s instead if IntDomain", subexpr_val.domain_kind))
		}
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: subexpr_val.Get_int().Neg()}
	case *imp.AddExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Add expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Add(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.SubExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Sub expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Sub(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.MulExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Mul expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Mul(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.DivExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Div expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Div(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.ModExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && rhs_val.domain_kind == IntDomainKind) {
			write_error(node_location, "Add expected LHS and RHS to be integer domain values, but are not")
		}
		result_intdom := lhs_val.Get_int().Mod(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: result_intdom}
	case *imp.ParenExpr:
		return interpreter.Eval_expr(node_location, expr_ty.Subexpr, abs_mem)
	case *imp.ArrayIndexExpr:
		arr_val := interpreter.Eval_expr(node_location, expr_ty.Base, abs_mem)
		index_val := interpreter.Eval_expr(node_location, expr_ty.Index, abs_mem)
		if arr_val.domain_kind != ArrayDomainKind {
			write_error(node_location, fmt.Sprintf("'%s' : expected arr to have arr domain type", expr_ty))
		}
		result_val := arr_val.Get_array().Index(index_val.Get_int())
		return result_val
	case *imp.EqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if lhs_val.domain_kind != rhs_val.domain_kind {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different", expr_ty))
		}
		result_val := lhs_val.Get_int().Eq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.NeqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if lhs_val.domain_kind != rhs_val.domain_kind {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different", expr_ty))
		}
		result_val := lhs_val.Get_int().Neq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	case *imp.LeqExpr:
		lhs_val := interpreter.Eval_expr(node_location, expr_ty.Lhs, abs_mem)
		rhs_val := interpreter.Eval_expr(node_location, expr_ty.Rhs, abs_mem)
		if !(lhs_val.domain_kind == IntDomainKind && lhs_val.domain_kind == rhs_val.domain_kind) {
			write_error(node_location, fmt.Sprintf("'%s' : types of LHS and RHS are different (%s vs %s)", expr_ty, lhs_val.domain_kind, rhs_val.domain_kind))
		}
		result_val := lhs_val.Get_int().Leq(rhs_val.Get_int())
		return AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: BoolDomainKind, bool_domain: result_val}
	}
	return AbstractValue[IntDomainImpl, ArrayDomainImpl]{}
}

func get_varname_from_lvalue(expr imp.Expr) string {
	switch expr_ty := expr.(type) {
	case *imp.VarExpr:
		return expr_ty.Name
	}
	panic(fmt.Sprintf("get_varname_from_lvalue: unimplemented expr type %T", expr))
}

func (interpreter *ImpFunctionInterpreter[IntDomainImpl, ArrayDomainImpl]) Step(in_state AbstractState[IntDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, ArrayDomainImpl] {
	cfg_node, cfg_node_exists := interpreter.func_cfg_map[interpreter.func_name].Node_map[in_state.node_location.Id]
	if !cfg_node_exists {
		write_error(create_empty_node_location(), fmt.Sprintf("The designated CFG Node %s doesn't exist", in_state.node_location))
	}

	// When we receive a new pair, update the global state with its join
	global_state, _ := interpreter.abstract_mem.pre_mem[in_state.node_location.Id]
	state_changed := global_state.Join_inplace(in_state.abstract_mem)
	if !state_changed && interpreter.abstract_mem.n_visits[in_state.node_location.Id] > 0 {
		// no updates to the state
		write_info(in_state.node_location, "No updates to node state")
		return nil
	}
	write_update_node(in_state.node_location, global_state.String())
	interpreter.abstract_mem.n_visits[in_state.node_location.Id]++
	// Executed on the joined node state
	in_state.abstract_mem = global_state.Clone()

	var return_states []AbstractState[IntDomainImpl, ArrayDomainImpl]
	switch cfg_node := cfg_node.(type) {
	case *CFGNode:
		switch stmt := cfg_node.Ast.(type) {
		case *imp.AssignStmt:
			// assignment should overwrite the value, instead of join
			rhs_val := interpreter.Eval_expr(in_state.node_location, stmt.Rhs, in_state.abstract_mem)
			switch lhs_ty := stmt.Lhs.(type) {
			case *imp.VarExpr:
				_, var_exists := in_state.abstract_mem[lhs_ty.Name]
				if var_exists {
					if in_state.abstract_mem[lhs_ty.Name].domain_kind != rhs_val.domain_kind {
						write_error(in_state.node_location, "LHS and RHS domain type does not match")
						return nil
					}
				}
				in_state.abstract_mem[lhs_ty.Name] = rhs_val
			}

		case *imp.SkipStmt:
			// do nothing
		default:
			panic("unimplemented")
		}
	case *CFGCondNode:
		cond_edge, is_cond_edge := interpreter.func_cfg_map[interpreter.func_name].Edge_map_from[in_state.node_location.Id].(*CFGCondEdge)
		if !is_cond_edge {
			write_error(in_state.node_location, "Condition stmt does not have outgoing edge of CondEdge type.")
		}
		// If at in_state the prop evaluates to either true or false,
		// We can just execute only the corresponding branch.
		// Otherwise filter for each branch and join the result
		cond_val := interpreter.Eval_expr(in_state.node_location, cfg_node.Cond_expr, in_state.abstract_mem)

		// Try to represent it as SimpleProp
		cond_simpleprop, simpleprop_success := algebra.Imp_expr_to_simple_prop(cfg_node.Cond_expr)
		if !simpleprop_success {
			write_warning(in_state.node_location, fmt.Sprintf("Could not represent '%s' as SimpleProp. Analysis precision may severely deterioriate.", cfg_node.Cond_expr))
		}
		if cond_val.Get_bool().IsBot() { // dead branch
			return nil
		}

		if (cond_val.Get_bool().IsTrue() || cond_val.Get_bool().IsTop()) && cond_edge.To_true_node_loc.NodeExists() {
			// run just the true branch on filter_true(in_state)
			new_state := in_state.Clone()
			if simpleprop_success {
				true_filters := domain.Filter_true_query_simpleprop(cond_simpleprop)
				for _, filter := range true_filters {
					lhs_name := get_varname_from_lvalue(filter.Term_expr)
					rhs_dom_val := interpreter.Eval_expr(in_state.node_location, filter.Rhs_expr, in_state.abstract_mem)
					if rhs_dom_val.domain_kind == IntDomainKind {
						updated_intdom := new_state.abstract_mem[lhs_name].Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
						new_state.abstract_mem[lhs_name] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
					}
				}
			}
			return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_true_node_loc, abstract_mem: new_state.abstract_mem})
		}

		if (cond_val.Get_bool().IsFalse() || cond_val.Get_bool().IsTop()) && cond_edge.To_false_node_loc.NodeExists() {
			// run just the false branch on filter_false(in_state)
			new_state := in_state.Clone()
			if simpleprop_success {
				false_filters := domain.Filter_false_query_simpleprop(cond_simpleprop)
				for _, filter := range false_filters {
					lhs_name := get_varname_from_lvalue(filter.Term_expr)
					rhs_dom_val := interpreter.Eval_expr(in_state.node_location, filter.Rhs_expr, in_state.abstract_mem)
					if rhs_dom_val.domain_kind == IntDomainKind {
						updated_intdom := new_state.abstract_mem[lhs_name].Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
						new_state.abstract_mem[lhs_name] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
					}
				}
			}

			return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_false_node_loc, abstract_mem: new_state.abstract_mem})
		}

		if cond_val.Get_bool().IsTop() {
			// run both branches
			if cond_edge.To_true_node_loc.NodeExists() { // true stmt exists
				new_state := in_state.Clone()
				if simpleprop_success { // apply filter
					true_filters := domain.Filter_true_query_simpleprop(cond_simpleprop)
					for _, filter := range true_filters {
						lhs_name := get_varname_from_lvalue(filter.Term_expr)
						rhs_dom_val := interpreter.Eval_expr(in_state.node_location, filter.Rhs_expr, in_state.abstract_mem)
						if rhs_dom_val.domain_kind == IntDomainKind {
							updated_intdom := new_state.abstract_mem[lhs_name].Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
							new_state.abstract_mem[lhs_name] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
						}
					}
				}
				return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_true_node_loc, abstract_mem: new_state.abstract_mem})
			}

			if cond_edge.To_false_node_loc.NodeExists() { // false stmt exists
				new_state := in_state.Clone()
				if simpleprop_success { // apply filter
					false_filters := domain.Filter_false_query_simpleprop(cond_simpleprop)
					for _, filter := range false_filters {
						lhs_name := get_varname_from_lvalue(filter.Term_expr)
						rhs_dom_val := interpreter.Eval_expr(in_state.node_location, filter.Rhs_expr, in_state.abstract_mem)
						if rhs_dom_val.domain_kind == IntDomainKind {
							updated_intdom := new_state.abstract_mem[lhs_name].Get_int().Filter(filter.Query_type, rhs_dom_val.Get_int())
							new_state.abstract_mem[lhs_name] = AbstractValue[IntDomainImpl, ArrayDomainImpl]{domain_kind: IntDomainKind, int_domain: updated_intdom}
						}
					}
				}
				return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: cond_edge.To_false_node_loc, abstract_mem: new_state.abstract_mem})
			}
		}
	}

	switch outgoing_edge := interpreter.func_cfg_map[interpreter.func_name].Edge_map_from[in_state.node_location.Id].(type) {
	case *CFGEdge:
		new_state := in_state.Clone()
		return_states = append(return_states, AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: outgoing_edge.To_node_loc, abstract_mem: new_state.abstract_mem})
	case *CFGCondEdge:
		// handle condition edges within their respective stmt handling
	}
	return return_states
}
