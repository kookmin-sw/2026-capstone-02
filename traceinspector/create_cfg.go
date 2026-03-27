package traceinspector

import (
	"fmt"
	"go/ast"
	"go/token"
)

func Print_cfg(file *ast.File, fset *token.FileSet) {
	for _, decls := range file.Decls {
		switch decl_node := decls.(type) {
		case *ast.FuncDecl:
			if decl_node.Name.Name == "main" {
				cfg_graph := CFGGraph{}
				cfg_creator := CFGGraphCreator{fset: fset, cfg_graph: &cfg_graph, next_node_index: 1}
				cfg_creator.create_cfg_method(decl_node.Body)
				fmt.Println(string(cfg_creator.cfg_graph.to_json()))
			}
		}
	}
}
