package traceinspector

import (
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"traceinspector/imp"
)

// if node id is leq 0, then the node doesn't exist
type CFGGraphCreator struct {
	func_name         imp.ImpFunctionName // name of the function
	fset              *token.FileSet
	Cfg_graph         *CFGGraph
	next_node_id      NodeID       // the next available node id
	next_edge_id      EdgeID       // the next available edge id
	cfg_context_stack []CFGContext // stack holding the graph context

}

// Map from function names to CFGGraph
type FunctionCFGMap map[imp.ImpFunctionName]*CFGGraph

type CFGContext interface {
	isCFGContext()
}

type CFGLoopContext struct {
	head_node_loc CFGNodeLocation // loc of the loop head(condition node)
	exit_node_loc CFGNodeLocation // loc of the node after the loop
}

func (CFGLoopContext) isCFGContext() {}

type CFGBranchContext struct {
	exit_node_loc CFGNodeLocation // loc of the node after the branch(join node)
}

func (CFGBranchContext) isCFGContext() {}

// Return the topmost loop context
func (creator *CFGGraphCreator) get_top_loop_context() *CFGLoopContext {
	for stack_index := len(creator.cfg_context_stack) - 1; stack_index >= 0; stack_index-- {
		loop_context, is_loop_context := creator.cfg_context_stack[stack_index].(CFGLoopContext)
		if is_loop_context {
			return &loop_context
		}
	}
	return nil
}

// Return the next stmt node ID to evaluate(link to), given the current get_top_context_destination
// If in a branch, it's the stmt after the branch
// If in a loop, it's the loop head(condition node)
func (creator *CFGGraphCreator) get_top_context_destination() CFGNodeLocation {
	for stack_index := len(creator.cfg_context_stack) - 1; stack_index >= 0; stack_index-- {
		switch ctx := creator.cfg_context_stack[stack_index].(type) {
		case CFGLoopContext:
			if ctx.exit_node_loc.Id > 0 {
				return ctx.head_node_loc
			}
		case CFGBranchContext:
			if ctx.exit_node_loc.Id > 0 {
				return ctx.exit_node_loc
			}
		}
	}
	return CFGNodeLocation{creator.func_name, 0, 0}
}

func (creator *CFGGraphCreator) push_branch_context(cond_node_loc CFGNodeLocation, exit_node_loc CFGNodeLocation) {
	creator.cfg_context_stack = append(creator.cfg_context_stack, CFGBranchContext{exit_node_loc})
}

func (creator *CFGGraphCreator) push_loop_context(cond_node_loc CFGNodeLocation, exit_node_loc CFGNodeLocation) {
	creator.cfg_context_stack = append(creator.cfg_context_stack, CFGLoopContext{cond_node_loc, exit_node_loc})
}

func (creator *CFGGraphCreator) pop_context() {
	creator.cfg_context_stack = creator.cfg_context_stack[:len(creator.cfg_context_stack)-1]
}

func (graphcreator *CFGGraphCreator) create_cfg_node(imp_ast imp.Stmt, line_num int) CFGNodeLocation {
	current_node_index := graphcreator.next_node_id
	loc := CFGNodeLocation{graphcreator.func_name, current_node_index, line_num}
	graphcreator.Cfg_graph.Node_map[current_node_index] = &CFGNode{Ast: imp_ast, Loc: loc, Code: fmt.Sprintf("%s", imp_ast), Node_type: node_basic, Line_num: line_num}
	graphcreator.next_node_id++
	return loc
}

func (graphcreator *CFGGraphCreator) create_cfg_cond_node(imp_ast imp.Expr, line_num int, is_loop_head bool) CFGNodeLocation {
	current_node_index := graphcreator.next_node_id
	loc := CFGNodeLocation{graphcreator.func_name, current_node_index, line_num}
	graphcreator.Cfg_graph.Node_map[current_node_index] = &CFGCondNode{Cond_expr: imp_ast, Loc: loc, Code: fmt.Sprintf("%s", imp_ast), Node_type: node_cond, Line_num: line_num, Is_loop_head: is_loop_head}
	graphcreator.next_node_id++
	return loc
}

func (graphcreator *CFGGraphCreator) create_cfg_edge(from_loc CFGNodeLocation, to_loc CFGNodeLocation, label string) {
	if from_loc.Id > 0 && to_loc.Id > 0 {
		edge := CFGEdge{Loc: CFGEdgeLocation{graphcreator.func_name, graphcreator.next_edge_id}, From_node_loc: from_loc, To_node_loc: to_loc, Label: label}
		graphcreator.Cfg_graph.Edge_map_from[from_loc.Id] = &edge
		graphcreator.Cfg_graph.Edge_map_to[to_loc.Id] = append(graphcreator.Cfg_graph.Edge_map_to[to_loc.Id], &edge)
		graphcreator.next_edge_id++
	}
}

func (graphcreator *CFGGraphCreator) create_cfg_cond_edge(from_loc CFGNodeLocation, to_true_loc CFGNodeLocation, to_false_loc CFGNodeLocation) {
	if from_loc.Id > 0 && (to_true_loc.Id > 0 || to_false_loc.Id > 0) {
		edge := CFGCondEdge{Loc: CFGEdgeLocation{graphcreator.func_name, graphcreator.next_edge_id}, From_node_loc: from_loc, To_true_node_loc: to_true_loc, To_false_node_loc: to_false_loc}
		graphcreator.Cfg_graph.Edge_map_from[from_loc.Id] = &edge
		if to_true_loc.Id > 0 {
			graphcreator.Cfg_graph.Edge_map_to[to_true_loc.Id] = append(graphcreator.Cfg_graph.Edge_map_to[to_true_loc.Id], &edge)
		}
		if to_false_loc.Id > 0 {
			graphcreator.Cfg_graph.Edge_map_to[to_false_loc.Id] = append(graphcreator.Cfg_graph.Edge_map_to[to_false_loc.Id], &edge)
		}
		graphcreator.next_edge_id++
	}
}

// The driver function for creating the CFG graph. stmt is the current statement node.
// Returns the CFGNodeLocation of the created node.
func (graphcreator *CFGGraphCreator) create_cfg_method(stmts []imp.Stmt) CFGNodeLocation {
	if len(stmts) == 0 {
		return CFGNodeLocation{graphcreator.func_name, 0, 0}
	}
	next_node_loc := graphcreator.create_cfg_method(stmts[1:]) // slice[1:] returns empty slice for len 1 slice
	if next_node_loc.Id == 0 {
		// If there's no remaining statement, the next destination depends on context
		next_node_loc = graphcreator.get_top_context_destination()
	}
	var created_node_loc CFGNodeLocation = CFGNodeLocation{graphcreator.func_name, 0, 0}
	switch stmt_ty := stmts[0].(type) {
	case *imp.IfElseStmt:
		cond_node_id := graphcreator.create_cfg_cond_node(stmt_ty.Cond, stmt_ty.GetLineNum(), false)

		graphcreator.push_branch_context(cond_node_id, next_node_loc)

		// the node location of the starting node in true stmt flow
		true_node_loc := graphcreator.create_cfg_method(stmt_ty.True_stmt)
		if true_node_loc.Id == 0 {
			// true stmt empty, next destination is context-dependent
			true_node_loc = next_node_loc
		}

		// the node location of the starting node in the false stmt flow
		false_node_loc := graphcreator.create_cfg_method(stmt_ty.False_stmt)
		if false_node_loc.Id == 0 {
			// false stmt empty, next destination is context dependent
			false_node_loc = next_node_loc
		}

		// create edges from cond to true/false start node
		// graphcreator.create_cfg_edge(cond_node_id, true_node_id, "True")
		// graphcreator.create_cfg_edge(cond_node_id, false_node_id, "False")
		graphcreator.create_cfg_cond_edge(cond_node_id, true_node_loc, false_node_loc)

		graphcreator.pop_context()

		created_node_loc = cond_node_id

	case *imp.WhileStmt:
		cond_node_loc := graphcreator.create_cfg_cond_node(stmt_ty.Cond, stmt_ty.GetLineNum(), true)

		graphcreator.push_loop_context(cond_node_loc, next_node_loc)
		body_node_id := graphcreator.create_cfg_method(stmt_ty.Body_stmt)
		// graphcreator.create_cfg_edge(cond_node_id, body_node_id, "True")
		// graphcreator.create_cfg_edge(cond_node_id, next_node_id, "False")
		graphcreator.create_cfg_cond_edge(cond_node_loc, body_node_id, next_node_loc)
		graphcreator.pop_context()

		created_node_loc = cond_node_loc

	case *imp.BreakStmt:
		created_node_loc = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		ctx := graphcreator.get_top_loop_context()
		// link to loop exit
		graphcreator.create_cfg_edge(created_node_loc, ctx.exit_node_loc, "")

	case *imp.ContinueStmt:
		created_node_loc = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		ctx := graphcreator.get_top_loop_context()
		// link to loop head
		graphcreator.create_cfg_edge(created_node_loc, ctx.head_node_loc, "")

	case *imp.ReturnStmt:
		created_node_loc = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		// finish generation
	default:
		created_node_loc = graphcreator.create_cfg_node(stmts[0], stmt_ty.GetLineNum())
		graphcreator.create_cfg_edge(created_node_loc, next_node_loc, "")

	}
	return created_node_loc
}

// create and print the cfg into json
func Create_cfg(functions imp.ImpFunctionMap) FunctionCFGMap {
	var func_cfg_map FunctionCFGMap = make(FunctionCFGMap)
	for fun_name, fun := range functions {
		func_cfg_map[fun_name] = &CFGGraph{Node_map: make(map[NodeID]CFGNodeClass), Edge_map_from: map[NodeID]CFGEdgeClass{}, Edge_map_to: map[NodeID][]CFGEdgeClass{}}
		cfg_creator := CFGGraphCreator{func_name: fun_name, Cfg_graph: func_cfg_map[fun_name], next_node_id: 1}
		entry_node_loc := cfg_creator.create_cfg_method(fun.Body)
		func_cfg_map[fun_name].Entry_node = entry_node_loc
	}
	return func_cfg_map
}

func Print_cfg_map_json(cfgs FunctionCFGMap) {
	// result, _ := json.Marshal(func_cfg_map)
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	enc.Encode(cfgs)
}

func Print_mermaid(cfg *CFGGraph) {
	// fmt.Println("```")
	fmt.Println(cfg.To_mermaid())
	// fmt.Println("```")
}
