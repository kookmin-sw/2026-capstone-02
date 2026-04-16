package traceinspector

import (
	"fmt"
	"traceinspector/imp"
)

type node_types string

type CFGNodeClass interface {
	is_CFGNodeClass()
	To_mermaid() string
}

const (
	node_basic node_types = "basic"
	node_cond  node_types = "cond"
)

type CFGId struct {
	Function_name string
	Id            int
}

type CFGNode struct {
	Ast       *imp.Stmt `json:"-"`
	Id        CFGId
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGCondNode struct {
	Ast       *imp.Expr `json:"-"`
	Id        CFGId
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGEdge struct {
	Id           CFGId
	From_node_id int
	To_node_id   int
	Label        string
}

func (node *CFGNode) is_CFGNodeClass() {}

func (node *CFGNode) To_mermaid() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id.Id, node.Code)
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id.Id, node.Code)
	}
	return ""
}

func (node *CFGCondNode) is_CFGNodeClass() {}

func (node *CFGCondNode) To_mermaid() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id.Id, node.Code)
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id.Id, node.Code)
	}
	return ""
}

func (node *CFGEdge) To_mermaid() string {
	if node.Label == "" {
		return fmt.Sprintf("%d --> %d", node.From_node_id, node.To_node_id)
	} else {
		return fmt.Sprintf("%d -- %s --> %d", node.From_node_id, node.Label, node.To_node_id)
	}
}

type CFGGraph struct {
	Nodes []CFGNodeClass
	Edges []CFGEdge
}
