package traceinspector

import (
	"fmt"
	"traceinspector/imp"
)

type node_types string

type CFGNodeClass interface {
	to_mermaind() string
}

const (
	node_basic node_types = "basic"
	node_cond  node_types = "cond"
)

type CFGNode struct {
	ast       *imp.Stmt `json:"-"`
	Id        int
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGCondNode struct {
	ast       *imp.Expr `json:"-"`
	Id        int
	Code      string
	Node_type node_types
	Line_num  int
}

type CFGEdge struct {
	Id           int
	From_node_id int
	To_node_id   int
	Label        string
}

func (node *CFGNode) to_mermaind() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id, node.Code)
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id, node.Code)
	}
	return ""
}

func (node *CFGCondNode) to_mermaind() string {
	switch node.Node_type {
	case node_basic:
		return fmt.Sprintf("%d[\"`%s`\"]", node.Id, node.Code)
	case node_cond:
		return fmt.Sprintf("%d{\"`%s`\"}", node.Id, node.Code)
	}
	return ""
}

func (node *CFGEdge) to_mermaind() string {
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
