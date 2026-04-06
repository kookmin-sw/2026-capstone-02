package traceinspector

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"traceinspector/imp"
)

// if node id is leq 0, then the node doesn't exist
type CFGGraphCreator struct {
	fset              *token.FileSet
	cfg_graph         *CFGGraph
	next_node_index   int          // the next available node id
	next_edge_index   int          // the next available edge id
	cfg_context_stack []CFGContext // stack holding the graph context

}

type CFGContext interface {
	isCFGContext()
}

type CFGLoopContext struct {
	head_node_id int // node ID of the loop head(condition node)
	exit_node_id int // node ID of the node after the loop
}

func (CFGLoopContext) isCFGContext() {}

type CFGBranchContext struct {
	cond_node_id int // node ID of the condition
	exit_node_id int // node ID of the node after the branch(join node)
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

func (creator *CFGGraphCreator) get_top_branch_context() *CFGBranchContext {
	for stack_index := len(creator.cfg_context_stack) - 1; stack_index >= 0; stack_index-- {
		branch_context, is_branch_context := creator.cfg_context_stack[stack_index].(CFGBranchContext)
		if is_branch_context {
			return &branch_context
		}
	}
	return nil
}

func (creator *CFGGraphCreator) push_branch_context(cond_node_id int, exit_node_id int) {
	creator.cfg_context_stack = append(creator.cfg_context_stack, CFGBranchContext{cond_node_id, exit_node_id})
}

func (creator *CFGGraphCreator) push_context(context CFGContext) {
	creator.cfg_context_stack = append(creator.cfg_context_stack, context)
}

func (creator *CFGGraphCreator) pop_context() {
	creator.cfg_context_stack = creator.cfg_context_stack[:len(creator.cfg_context_stack)-1]
}

func (graphcreator *CFGGraphCreator) create_cfg_node(code_string string, node_type node_types, line_num int) int {
	current_node_index := graphcreator.next_node_index
	graphcreator.cfg_graph.Nodes = append(graphcreator.cfg_graph.Nodes, CFGNode{Id: current_node_index, Code: code_string, Node_type: node_type, Line_num: line_num})
	graphcreator.next_node_index++
	return current_node_index
}

func (graphchreator *CFGGraphCreator) create_cfg_edge(from_id int, to_id int, label string) {
	if from_id > 0 && to_id > 0 {
		graphchreator.cfg_graph.Edges = append(graphchreator.cfg_graph.Edges, CFGEdge{Id: graphchreator.next_edge_index, From_node_id: from_id, To_node_id: to_id, Label: label})
		graphchreator.next_edge_index++
	}
}

// The driver function for creating the CFG graph. stmt is the current statement node.
// linkback, if not 0, equals the node id that an edge should be created from the current node to the linkback ID
func (graphcreator *CFGGraphCreator) create_cfg_method(stmts []imp.Stmt) int {
	if len(stmts) == 0 {
		return 0
	}
	exit_node_id := graphcreator.create_cfg_method(stmts[1:]) // slice[1:] returns empty slice for len 1 slice
	switch stmt_ty := stmts[0].(type) {
	case *imp.IfElseStmt:
		cond_node_id := graphcreator.create_cfg_node(fmt.Sprintf("%s", stmt_ty.Cond), node_cond, stmt_ty.Line_num)

		graphcreator.push_branch_context(cond_node_id, exit_node_id)
		true_node_id := graphcreator.create_cfg_method(stmt_ty.True_stmt)
		false_node_id := graphcreator.create_cfg_method(stmt_ty.False_stmt)
		graphcreator.create_cfg_edge(cond_node_id, true_node_id, "True")
		graphcreator.create_cfg_edge(cond_node_id, false_node_id, "False")
		graphcreator.pop_context()

		return cond_node_id
	case *imp.WhileStmt:
		cond_node_id := graphcreator.create_cfg_node(fmt.Sprintf("%s", stmt_ty.Cond), node_cond, stmt_ty.Line_num)

		graphcreator.push_branch_context(cond_node_id, exit_node_id)
		body_node_id := graphcreator.create_cfg_method(stmt_ty.Body_stmt)
		graphcreator.create_cfg_edge(cond_node_id, body_node_id, "True")
		graphcreator.create_cfg_edge(cond_node_id, exit_node_id, "False")
		graphcreator.pop_context()
	case *imp.BreakStmt:
		node_id := graphcreator.create_cfg_node("break", node_basic, stmt_ty.Line_num)
		ctx := graphcreator.get_top_loop_context()
		graphcreator.create_cfg_edge(node_id, ctx.exit_node_id, "")
	case *imp.ContinueStmt:
		node_id := graphcreator.create_cfg_node("continue", node_basic, stmt_ty.Line_num)
		ctx := graphcreator.get_top_loop_context()
		graphcreator.create_cfg_edge(node_id, ctx.head_node_id, "")
	case *imp.ReturnStmt:
		node_id := graphcreator.create_cfg_node(fmt.Sprintf("%s", stmt_ty), node_basic, stmt_ty.Line_num)
		return node_id

	}
}

// func (graphcreator *CFGGraphCreator) create_cfg_method(stmts []imp.Stmt) int {
// 	if stmts == nil {
// 		return 0
// 	}
// 	for _, stmt := range stmts {
// 		switch stmt_ty := stmt.(type) {
// 		case *imp.IfElseStmt:
// 			if_node_id := graphcreator.create_cfg_node(stmt_ty.Cond, node_cond)
// 		}
// 	}

// 	switch node := node.(type) {
// 	case *imp:
// 		if len(node.List) > 0 {
// 			for _, subnode := range node.List {
// 				graphcreator.create_cfg_method(subnode)
// 			}
// 		}
// 		return graphcreator.prev_node_index
// 	case *ast.IfStmt:
// 		if_node_id := graphcreator.create_cfg_node(node.Cond, node_cond)
// 		graphcreator.create_cfg_edge(if_node_id, "")
// 		graphcreator.prev_node_index = if_node_id
// 		if node.Body != nil {
// 			// true body
// 			for index, body_node := range node.Body.List {
// 				if index == 0 {
// 					body_node_id := graphcreator.create_cfg_node(body_node, node_basic)
// 					graphcreator.create_cfg_edge(body_node_id, "true")
// 					graphcreator.prev_node_index = body_node_id
// 				} else {
// 					graphcreator.prev_node_index = graphcreator.create_cfg_method(body_node)
// 				}
// 				graphcreator.prev_node_index_if_else = graphcreator.prev_node_index
// 			}
// 		}
// 		if node.Else != nil {
// 			switch else_node := node.Else.(type) {
// 			case *ast.BlockStmt:
// 				// else body
// 				for index, body_node := range else_node.List {
// 					if index == 0 {
// 						body_node_id := graphcreator.create_cfg_node(body_node, node_basic)
// 						graphcreator.create_cfg_edge(body_node_id, "true")
// 						graphcreator.prev_node_index = body_node_id
// 					} else {
// 						graphcreator.prev_node_index = graphcreator.create_cfg_method(body_node)
// 					}
// 				}
// 			default:
// 				graphcreator.prev_node_index = if_node_id
// 				else_node_id := graphcreator.create_cfg_node(node.Else, node_basic)
// 				graphcreator.create_cfg_edge(else_node_id, "false")
// 			}
// 		}
// 		return graphcreator.prev_node_index
// 	default:
// 		node_id := graphcreator.create_cfg_node(n, node_basic)
// 		graphcreator.create_cfg_edge(node_id, "")
// 		graphcreator.prev_node_index = node_id
// 		return node_id
// 	}
// }

// create and print the cfg into json
func Print_cfg(file *ast.File, fset *token.FileSet) {
	var func_cfg_map map[string]*CFGGraph = make(map[string]*CFGGraph)
	for _, decls := range file.Decls {
		switch decl_node := decls.(type) {
		case *ast.FuncDecl:
			{
				func_cfg_map[decl_node.Name.Name] = &CFGGraph{}
				cfg_creator := CFGGraphCreator{fset: fset, cfg_graph: func_cfg_map[decl_node.Name.Name], next_node_index: 1}
				cfg_creator.create_cfg_method(decl_node.Body)
			}
		}
	}
	// result, _ := json.Marshal(func_cfg_map)
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	// result, _ := json.MarshalIndent(func_cfg_map, "", "    ") // return formatted
	// fmt.Println(string(result))
	enc.Encode(func_cfg_map)
}
