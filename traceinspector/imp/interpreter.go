package imp

import (
	"fmt"
	"strings"
)

type ImpState struct {
	vars                  map[string]ImpValues // all exprs are reduced to go values
	current_function_name string
	return_value          ImpValues // The return value of the current local scope, if exists
}

// This designates the control flow result of executing statements
type ControlflowResult int

const (
	ControlNormal ControlflowResult = iota
	ControlBreak
	ControlContinue
	ControlReturn
)

type ImpInterpreter struct {
	States    []*ImpState // a stack of program states to represent scopes
	Functions map[string]ImpFunction
}

// Interpret Imp code starting from main()
func (interpreter *ImpInterpreter) Interpret_main() {
	interpreter.eval_function_call("main", nil, 0)
}

func (interpreter *ImpInterpreter) get_top_state() *ImpState {
	return interpreter.States[len(interpreter.States)-1]
}

// Return the topmost variable name and bool indicating if the variable exists
func (interpreter *ImpInterpreter) get_variable(name string) (ImpValues, bool) {
	toplevel_func_name := interpreter.States[len(interpreter.States)-1].current_function_name
	for stack_index := len(interpreter.States) - 1; stack_index >= 0; stack_index-- {
		// Check that the scope is within the function stack
		if toplevel_func_name != interpreter.States[stack_index].current_function_name {
			break
		}

		var_value, var_exists := interpreter.States[stack_index].vars[name]
		if var_exists {
			return var_value, true
		}
	}
	return nil, false
}

func (interpreter *ImpInterpreter) push_state(state ImpState) {
	interpreter.States = append(interpreter.States, &state)
}

func (interpreter *ImpInterpreter) pop_state() {
	interpreter.States = interpreter.States[:len(interpreter.States)-1]
}

func (ImpInterpreter *ImpInterpreter) deepcopy_impvalue(val ImpValues) ImpValues {
	switch val_ty := val.(type) {
	case *IntVal:
		return &IntVal{val: val_ty.val}
	case *BoolVal:
		return &BoolVal{val: val_ty.val}
	case *StringVal:
		return &StringVal{val: val_ty.val}
	case *ArrayVal:
		arrayval := ArrayVal{element_type: val_ty.element_type}
		var copied_slice []ImpValues
		for _, elem := range val_ty.val {
			copied_slice = append(copied_slice, ImpInterpreter.deepcopy_impvalue(elem))
		}
		arrayval.val = copied_slice
		return &arrayval
	default:
		panic(fmt.Sprintf("deepcopy_impvalue: Unknown value type %T\n", val))
	}
}

////////////////////////

func (interpreter *ImpInterpreter) eval_VarExpr(node VarExpr) ImpValues {
	var_value, var_exists := interpreter.get_variable(node.name)
	if !var_exists {
		panic(fmt.Sprintf("Line %d: Unknown variable '%s'", node.Node.Line_num, node.name))
	}
	return var_value
}

func (interpreter *ImpInterpreter) eval_Expr_lvalue(lhs Expr, rhs_ty ImpTypes) ImpValues {
	lhs_var, lhs_is_var := lhs.(*VarExpr)
	if lhs_is_var {
		_, lhs_exists := interpreter.get_variable(lhs_var.name)
		if !lhs_exists {
			switch ty := rhs_ty.(type) {
			case IntType:
				interpreter.get_top_state().vars[lhs_var.name] = &IntVal{}
			case BoolType:
				interpreter.get_top_state().vars[lhs_var.name] = &BoolVal{}
			case ArrayType:
				interpreter.get_top_state().vars[lhs_var.name] = &ArrayVal{element_type: ty.Element_type}
			default:
				panic(fmt.Sprintf("Line %d: Unknown rhs type %T", lhs_var.Line_num, ty))
			}
		}
		var_val, _ := interpreter.get_variable(lhs_var.name)
		return var_val
	} else {
		return interpreter.eval_Expr(lhs)
	}
}

func (interpreter *ImpInterpreter) eval_IntLitExpr(node IntLitExpr) ImpValues {
	return &IntVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_BoolLitExpr(node BoolLitExpr) ImpValues {
	return &BoolVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_StringLitExpr(node StringLitExpr) ImpValues {
	return &StringVal{val: node.value}
}

func (interpreter *ImpInterpreter) eval_ArrayLitExpr(node ArrayLitExpr) ImpValues {
	var elem_vals []ImpValues
	var elem_type ImpTypes
	for _, elem := range node.elements {
		elem_val := interpreter.eval_Expr(elem)
		if elem_type == nil {
			elem_type = get_type(elem_val)
		}
		if !check_val_type_match(elem_val, elem_type) {
			panic(fmt.Sprintf("Line %d: Array element type is identified as %s, but got expr '%s' with value %s\n", node.Line_num, elem_type, elem, elem_val))
		}
		elem_vals = append(elem_vals, elem_val)
	}
	return &ArrayVal{element_type: elem_type, val: elem_vals}
}

func (interpreter *ImpInterpreter) eval_AddExpr(node AddExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of addition should be an int value, but got '%s'", node.Line_num, node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of addition should be an int value, but got '%s'", node.Line_num, node.rhs))
	}
	return &IntVal{val: lhs_val.val + rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_SubExpr(node SubExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of subtraction should be an int value, but got '%s'", node.Line_num, node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("RHS of subtraction should be an int value, but got '%s'", node.rhs))
	}
	return &IntVal{val: lhs_val.val - rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_MulExpr(node MulExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of multiplication should be an int value, but got '%s'", node.Line_num, node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of multiplication should be an int value, but got '%s'", node.Line_num, node.rhs))
	}
	return &IntVal{val: lhs_val.val * rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_DivExpr(node DivExpr) ImpValues {
	lhs_val, lhs_is_int := interpreter.eval_Expr(node.lhs).(*IntVal)
	rhs_val, rhs_is_int := interpreter.eval_Expr(node.rhs).(*IntVal)

	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: LHS of division should be an int value, but got '%s'", node.Line_num, node.lhs))
	}

	if !rhs_is_int {
		panic(fmt.Sprintf("Line %d: RHS of division should be an int value, but got '%s'", node.Line_num, node.rhs))
	}
	return &IntVal{val: lhs_val.val / rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_ParenExpr(node ParenExpr) ImpValues {
	return interpreter.eval_Expr(node.subexpr)
}

func (interpreter *ImpInterpreter) eval_ArrayIndexExpr(node ArrayIndexExpr) ImpValues {
	index_val, index_is_int := interpreter.eval_Expr(node.index).(*IntVal)
	if !index_is_int {
		panic(fmt.Sprintf("Line %d: Index of array indexing should be an int value, but got '%s'", node.Line_num, node.index))
	}
	base_val, base_is_arrayval := interpreter.eval_Expr(node.base).(*ArrayVal)
	if !base_is_arrayval {
		panic(fmt.Sprintf("Line %d: Expr '%s' is not an array", node.Line_num, node.base))
	}
	return base_val.val[index_val.val]
}

func (interpreter *ImpInterpreter) eval_EqExpr(node EqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	if !check_vals_type_equal(lhs_val, rhs_val) {
		panic(fmt.Sprintf("Line %d: Unsupported '==' between '%s' and '%s'", node.Line_num, lhs_val, rhs_val))
	}
	switch lhs_val := lhs_val.(type) {
	case *IntVal:
		rhs_val, _ := rhs_val.(*IntVal)
		return &BoolVal{val: lhs_val.val == rhs_val.val}
	case *BoolVal:
		rhs_val, _ := rhs_val.(*BoolVal)
		return &BoolVal{val: lhs_val.val == rhs_val.val}
	case *StringVal:
		rhs_val, _ := rhs_val.(*StringVal)
		return &BoolVal{val: lhs_val.val == rhs_val.val}
	case *NoneVal:
		return &BoolVal{val: true}
	default:
		panic(fmt.Sprintf("Line %d: Unsupported '==' between %s and %s", node.Line_num, lhs_val, rhs_val))
	}
}

func (interpreter *ImpInterpreter) eval_NeqExpr(node NeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	if !check_vals_type_equal(lhs_val, rhs_val) {
		panic(fmt.Sprintf("Line %d: Unsupported '!=' between %s and %s", node.Line_num, lhs_val, rhs_val))
	}
	switch lhs_val := lhs_val.(type) {
	case *IntVal:
		rhs_val, _ := rhs_val.(*IntVal)
		return &BoolVal{val: lhs_val.val != rhs_val.val}
	case *BoolVal:
		rhs_val, _ := rhs_val.(*BoolVal)
		return &BoolVal{val: lhs_val.val != rhs_val.val}
	case *StringVal:
		rhs_val, _ := rhs_val.(*StringVal)
		return &BoolVal{val: lhs_val.val != rhs_val.val}
	case *NoneVal:
		return &BoolVal{val: false}
	default:
		panic(fmt.Sprintf("Line %d: Unsupported '!=' between %s and %s", node.Line_num, lhs_val, rhs_val))
	}
}

func (interpreter *ImpInterpreter) eval_LessthanExpr(node LessthanExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Lessthan operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{val: lhs_intvar.val < rhs_intvar.val}
}

func (interpreter *ImpInterpreter) eval_GreaterthanExpr(node GreaterthanExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Greaterthan operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{val: lhs_intvar.val > rhs_intvar.val}
}

func (interpreter *ImpInterpreter) eval_LeqExpr(node LeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Leq operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{val: lhs_intvar.val <= rhs_intvar.val}
}

func (interpreter *ImpInterpreter) eval_GeqExpr(node GeqExpr) ImpValues {
	lhs_val := interpreter.eval_Expr(node.lhs)
	rhs_val := interpreter.eval_Expr(node.rhs)
	lhs_intvar, lhs_is_int := lhs_val.(*IntVal)
	rhs_intvar, rhs_is_int := rhs_val.(*IntVal)
	if !(lhs_is_int && rhs_is_int) {
		panic(fmt.Sprintf("Line %d: Geq operator must be applied between two integer values", node.Line_num))
	}
	return &BoolVal{val: lhs_intvar.val >= rhs_intvar.val}
}

func (interpreter *ImpInterpreter) eval_NegExpr(node NegExpr) ImpValues {
	subexpr_val, subexpr_is_int := interpreter.eval_Expr(node.subexpr).(*IntVal)
	if !subexpr_is_int {
		panic(fmt.Sprintf("Line %d: Subexpr %s of Unary neg operator should be of type int", node.Line_num, node.subexpr))
	}
	return &IntVal{val: -subexpr_val.val}
}

func (interpreter *ImpInterpreter) eval_NotExpr(node NotExpr) ImpValues {
	subexpr_val, subexpr_is_bool := interpreter.eval_Expr(node.subexpr).(*BoolVal)
	if !subexpr_is_bool {
		panic(fmt.Sprintf("Line %d: Subexpr %s of NOT operator should be of type bool", node.Line_num, node.subexpr))
	}
	return &BoolVal{val: !subexpr_val.val}
}

func (interpreter *ImpInterpreter) eval_AndExpr(node AndExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("Line %d: LHS of AND should be a bool value, but got '%s'", node.Line_num, node.lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("Line %d: RHS of AND should be a bool value, but got '%s'", node.Line_num, node.rhs))
	}
	return &BoolVal{val: lhs_val.val && rhs_val.val}
}

func (interpreter *ImpInterpreter) eval_OrExpr(node OrExpr) ImpValues {
	lhs_val, lhs_is_bool := interpreter.eval_Expr(node.lhs).(*BoolVal)
	rhs_val, rhs_is_bool := interpreter.eval_Expr(node.rhs).(*BoolVal)

	if !lhs_is_bool {
		panic(fmt.Sprintf("Line %d: LHS of OR should be a bool value, but got '%s'", node.Line_num, node.lhs))
	}

	if !rhs_is_bool {
		panic(fmt.Sprintf("Line %d: RHS of OR should be a bool value, but got '%s'", node.Line_num, node.rhs))
	}
	return &BoolVal{val: lhs_val.val || rhs_val.val}
}

// Imp is pass-by-value for int/bool, but arrays are passed references
func (interpreter *ImpInterpreter) eval_function_call(func_name string, args []Expr, line_num int) ImpValues {
	// copy values if primitive
	prepare_args := func(arg ImpValues) ImpValues {
		switch arg_ty := arg.(type) {
		case *IntVal:
			return &IntVal{val: arg_ty.val}
		case *BoolVal:
			return &BoolVal{val: arg_ty.val}
		case *ArrayVal:
			return arg
		}
		panic(fmt.Sprintf("Line %d: Unknown arg type '%s' for function %s\n", line_num, get_type(arg), func_name))
	}
	func_local_state := ImpState{vars: make(map[string]ImpValues), current_function_name: func_name, return_value: &NoneVal{}}
	imp_function, function_exists := interpreter.Functions[func_name]
	if !function_exists {
		panic(fmt.Sprintf("Line %d: Unknown function '%s'\n", line_num, func_name))
	}
	for index, arg_expr := range args {
		arg_info := imp_function.Arg_pairs[index]
		arg_val := prepare_args(interpreter.eval_Expr(arg_expr))
		if !check_val_type_match(arg_val, arg_info.arg_type) {
			panic(fmt.Sprintf("Line %d: Argument '%s' of function '%s' is defined as type %s, but passed expr '%s' of type %s", line_num, arg_info.name, func_name, arg_info.arg_type, arg_expr, get_type(arg_val)))
		}
		func_local_state.vars[arg_info.name] = arg_val
	}
	interpreter.push_state(func_local_state)
	interpreter.eval_Stmt(imp_function.Body)
	return_value := interpreter.get_top_state().return_value
	interpreter.pop_state()

	return return_value
}

func (interpreter *ImpInterpreter) eval_CallExpr(node CallExpr) ImpValues {
	return interpreter.eval_function_call(node.func_name, node.args, node.Line_num)
}

func (interpreter *ImpInterpreter) eval_MakeArrayExpr(node MakeArrayExpr) ImpValues {
	len_node := interpreter.eval_Expr(node.size)
	len_intval, len_is_int := len_node.(*IntVal)
	if !len_is_int {
		panic(fmt.Sprintf("Line %d: %s - length expression %s is not an integer value", node.Line_num, node, node.size))
	}
	default_val := interpreter.eval_Expr(node.value)
	generated := make([]ImpValues, len_intval.val)
	for i := 0; i < len_intval.val; i++ {
		generated[i] = interpreter.deepcopy_impvalue(default_val)
	}
	return &ArrayVal{get_type(default_val), generated}
}

func (interpreter *ImpInterpreter) eval_LenExpr(node LenExpr) ImpValues {
	array_node := interpreter.eval_Expr(node.subexpr)
	array_val, is_array := array_node.(*ArrayVal)
	if !is_array {
		panic(fmt.Sprintf("Line %d: len() - Non-array value %s passed to len()", node.Line_num, node.subexpr))
	}
	return &IntVal{val: len(array_val.val)}
}

func (interpreter *ImpInterpreter) eval_Expr(node Expr) ImpValues {
	switch node_ty := node.(type) {
	case *VarExpr:
		return interpreter.eval_VarExpr(*node_ty)
	case *IntLitExpr:
		return interpreter.eval_IntLitExpr(*node_ty)
	case *BoolLitExpr:
		return interpreter.eval_BoolLitExpr(*node_ty)
	case *StringLitExpr:
		return interpreter.eval_StringLitExpr(*node_ty)
	case *ArrayLitExpr:
		return interpreter.eval_ArrayLitExpr(*node_ty)
	case *AddExpr:
		return interpreter.eval_AddExpr(*node_ty)
	case *SubExpr:
		return interpreter.eval_SubExpr(*node_ty)
	case *MulExpr:
		return interpreter.eval_MulExpr(*node_ty)
	case *DivExpr:
		return interpreter.eval_DivExpr(*node_ty)
	case *ParenExpr:
		return interpreter.eval_ParenExpr(*node_ty)
	case *ArrayIndexExpr:
		return interpreter.eval_ArrayIndexExpr(*node_ty)
	case *EqExpr:
		return interpreter.eval_EqExpr(*node_ty)
	case *NeqExpr:
		return interpreter.eval_NeqExpr(*node_ty)
	case *LessthanExpr:
		return interpreter.eval_LessthanExpr(*node_ty)
	case *GreaterthanExpr:
		return interpreter.eval_GreaterthanExpr(*node_ty)
	case *LeqExpr:
		return interpreter.eval_LeqExpr(*node_ty)
	case *GeqExpr:
		return interpreter.eval_GeqExpr(*node_ty)
	case *NegExpr:
		return interpreter.eval_NegExpr(*node_ty)
	case *NotExpr:
		return interpreter.eval_NotExpr(*node_ty)
	case *AndExpr:
		return interpreter.eval_AndExpr(*node_ty)
	case *OrExpr:
		return interpreter.eval_OrExpr(*node_ty)
	case *CallExpr:
		return interpreter.eval_CallExpr(*node_ty)
	case *MakeArrayExpr:
		return interpreter.eval_MakeArrayExpr(*node_ty)
	case *LenExpr:
		return interpreter.eval_LenExpr(*node_ty)
	default:
		panic(fmt.Sprintf(" Unimplemented expr type %s", node))
	}
}

/////////////////////////////////
// statements

func (interpreter *ImpInterpreter) eval_SkipStmt(SkipStmt) ControlflowResult {
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_AssignStmt(node AssignStmt) ControlflowResult {
	rhs_val := interpreter.eval_Expr(node.rhs)
	switch lhs_loc := interpreter.eval_Expr_lvalue(node.lhs, get_type(rhs_val)).(type) {
	case *IntVal:
		rhs_intval, rhs_is_intval := rhs_val.(*IntVal)
		if !rhs_is_intval {
			panic(fmt.Sprintf("Line %d: Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.Line_num, node.rhs, get_type(rhs_val), node.lhs, get_type(lhs_loc)))
		}
		lhs_loc.val = rhs_intval.val
	case *BoolVal:
		rhs_boolval, rhs_is_boolval := rhs_val.(*BoolVal)
		if !rhs_is_boolval {
			panic(fmt.Sprintf("Line %d: Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.Line_num, node.rhs, get_type(rhs_val), node.lhs, get_type(lhs_loc)))
		}
		lhs_loc.val = rhs_boolval.val
	case *ArrayVal:
		rhs_arrval, rhs_is_arrayval := rhs_val.(*ArrayVal)
		if !rhs_is_arrayval {
			panic(fmt.Sprintf("Line %d: Attempted to assign RHS '%s' of type %s to LHS '%s' of type %s", node.Line_num, node.rhs, get_type(rhs_val), node.lhs, get_type(lhs_loc)))
		}
		lhs_loc.val = rhs_arrval.val
	default:
		panic(fmt.Sprintf("Line %d: LHS expr '%s' has unresolved value type %T\n", node.Line_num, node.lhs, lhs_loc))
	}
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_IfElseStmt(node IfElseStmt) ControlflowResult {
	cond_val := interpreter.eval_Expr(node.cond)
	cond_boolval, cond_is_bool := cond_val.(*BoolVal)
	if !cond_is_bool {
		panic(fmt.Sprintf("Line %d: If statement got non-boolean condition '%s'\n", node.Line_num, node.cond))
	}
	if cond_boolval.val {
		return interpreter.eval_Stmt(node.true_stmt)
	} else {
		return interpreter.eval_Stmt(node.false_stmt)
	}
}

func (interpreter *ImpInterpreter) eval_WhileStmt(node WhileStmt) ControlflowResult {
	for true {
		cond_val := interpreter.eval_Expr(node.cond)
		cond_boolval, cond_is_bool := cond_val.(*BoolVal)
		if !cond_is_bool {
			panic(fmt.Sprintf("Line %d: While statement got non-boolean condition '%s'\n", node.Line_num, node.cond))
		}
		if cond_boolval.val == false {
			break
		}
		stmt_result := interpreter.eval_Stmt(node.body_stmt)
		switch stmt_result {
		case ControlBreak:
			return ControlNormal
		case ControlContinue:
			continue
		case ControlReturn:
			return ControlReturn
		}
	}
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_BreakStmt(_ BreakStmt) ControlflowResult {
	return ControlBreak
}

func (interpreter *ImpInterpreter) eval_ContinueStmt(_ ContinueStmt) ControlflowResult {
	return ControlContinue
}

func (interpreter *ImpInterpreter) eval_IncStmt(node IncStmt) ControlflowResult {
	lhs_val_int, lhs_is_int := interpreter.eval_Expr(node.subexpr).(*IntVal)
	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: Attempted to increment non-integer value '%s'\n", node.Line_num, node))
	}
	lhs_val_int.val++
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_DecStmt(node DecStmt) ControlflowResult {
	lhs_val_int, lhs_is_int := interpreter.eval_Expr(node.subexpr).(*IntVal)
	if !lhs_is_int {
		panic(fmt.Sprintf("Line %d: Attempted to decrement non-integer value '%s'\n", node.Line_num, node))
	}
	lhs_val_int.val--
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_CallStmt(node CallStmt) ControlflowResult {
	interpreter.eval_function_call(node.func_name, node.args, node.Line_num)
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_PrintStmt(node PrintStmt) ControlflowResult {
	var outputs []string
	for _, arg := range node.args {
		arg_val := interpreter.eval_Expr(arg)
		outputs = append(outputs, fmt.Sprintf("%s", arg_val))
	}
	fmt.Print(strings.Join(outputs, " "))
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_ScanfStmt(node ScanfStmt) ControlflowResult {
	for index, fmt_str := range strings.Split(node.format_string, " ") {
		var imp_val ImpValues
		switch fmt_str {
		case "%d":
			var input int
			fmt.Scanf(fmt_str, &input)
			imp_val = &IntVal{val: input}
		case "%t":
			var input bool
			fmt.Scanf(fmt_str, &input)
			imp_val = &BoolVal{val: input}
		default:
			panic(fmt.Sprintf("Line %d: scanf - Unsupported formatting specifier %s\n", node.Line_num, fmt_str))
		}

		lhs_val := node.assign_locations[index]
		switch lhs_loc := interpreter.eval_Expr_lvalue(lhs_val, get_type(imp_val)).(type) {
		case *IntVal:
			rhs_intval, rhs_is_intval := imp_val.(*IntVal)
			if !rhs_is_intval {
				panic(fmt.Sprintf("Line %d: scanf - Attempted to assign input '%s' of type %T to variable '%s' of type %T\n", node.Line_num, imp_val, imp_val, lhs_val, lhs_loc))
			}
			lhs_loc.val = rhs_intval.val
		case *BoolVal:
			rhs_intval, rhs_is_boolval := imp_val.(*BoolVal)
			if !rhs_is_boolval {
				panic(fmt.Sprintf("Line %d: scanf - Attempted to assign input '%s' of type %T to variable '%s' of type %T\n", node.Line_num, imp_val, imp_val, lhs_val, lhs_loc))
			}
			lhs_loc.val = rhs_intval.val
		}
	}
	return ControlNormal
}

func (interpreter *ImpInterpreter) eval_ReturnStmt(node ReturnStmt) ControlflowResult {
	top_state := interpreter.get_top_state()
	top_state.return_value = interpreter.eval_Expr(node.arg)
	if top_state.current_function_name != "" {
		expected_return_type := interpreter.Functions[top_state.current_function_name].Return_type
		if !check_val_type_match(top_state.return_value, expected_return_type) {
			panic(fmt.Sprintf("Line %d: Function %s should return value of type %s, but actually returned '%s' of type %s\n", node.Line_num, top_state.current_function_name, expected_return_type, node.arg, top_state.return_value))
		}
	}
	return ControlReturn
}

// eval_Stmt evaluates a sequence of statements
// The bool return type designates whether the function has returned, and hence execution of the sequence should stop
func (interpreter *ImpInterpreter) eval_Stmt(nodes []Stmt) ControlflowResult {
	var returned ControlflowResult = ControlNormal
	for _, stmt := range nodes {
		switch stmt := stmt.(type) {
		case *SkipStmt:
			returned = interpreter.eval_SkipStmt(*stmt)
		case *AssignStmt:
			returned = interpreter.eval_AssignStmt(*stmt)
		case *IfElseStmt:
			returned = interpreter.eval_IfElseStmt(*stmt)
		case *WhileStmt:
			returned = interpreter.eval_WhileStmt(*stmt)
		case *BreakStmt:
			returned = interpreter.eval_BreakStmt(*stmt)
		case *ContinueStmt:
			returned = interpreter.eval_ContinueStmt(*stmt)
		case *IncStmt:
			returned = interpreter.eval_IncStmt(*stmt)
		case *DecStmt:
			returned = interpreter.eval_DecStmt(*stmt)
		case *CallStmt:
			returned = interpreter.eval_CallStmt(*stmt)
		case *PrintStmt:
			returned = interpreter.eval_PrintStmt(*stmt)
		case *ScanfStmt:
			returned = interpreter.eval_ScanfStmt(*stmt)
		case *ReturnStmt:
			returned = interpreter.eval_ReturnStmt(*stmt)
		}
		if returned != ControlNormal {
			return returned
		}
	}
	return returned
}
