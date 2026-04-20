package traceinspector

import (
	"fmt"
	"traceinspector/domain"
	"traceinspector/imp"
)

type AnalyzerSettings struct {
	loop_iters_before_widening int
}

// An AbstractState is the pair (l, M^#) ↪ (l', M^#') used in the abstract transition relation
// node_id: node ID to be interpreted
// abstract_mem: the input abstract memory state wrt the node should be interpreted
type AbstractState[IntDomainImpl domain.AbstractDomain[IntDomainImpl], BoolDomainImpl domain.AbstractDomain[BoolDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	node_id      NodeID
	abstract_mem AbstractNodeMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}

// Step: Given an input state (l, m^#), execute the abstract step relation for l under memory state m^#, and
// Return the subsequent states {(l', m^#')} ∈ P(L * M^#)
type AbstractSemantics[IntDomainImpl domain.AbstractDomain[IntDomainImpl], BoolDomainImpl domain.AbstractDomain[BoolDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] interface {
	Step(AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}

// Abstract transition semantics for Imp wrt to arbitrary abstract domain impelmentations

// ImpFunctionInterpreter performs abstract interpretation of a function body from a given initial state. The
// interpreter performs interpretation until it collects the fixpoint semantics for the function body, and hence the
// return value. The interpreter will spawn another ImpFunctionInterpreter in the case a function call is invoked.
type ImpFunctionInterpreter[IntDomainImpl domain.AbstractDomain[IntDomainImpl], BoolDomainImpl domain.AbstractDomain[BoolDomainImpl], ArrayDomainImpl domain.AbstractDomain[ArrayDomainImpl]] struct {
	func_cfg_map  FunctionCFGMap
	func_info_map imp.ImpFunctionMap
	abstract_mem  *FunctionAbstractMem[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]
}

func (interpreter *ImpFunctionInterpreter[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) Step(in_state AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl]) []AbstractState[IntDomainImpl, BoolDomainImpl, ArrayDomainImpl] {
	cfg_node, cfg_node_exists := interpreter.func_cfg_map.Node_map[in_state.node_id]
	if !cfg_node_exists {
		write_error(create_empty_node_location(), fmt.Sprintf("The designated CFG Node %d doesn't exist", in_state.node_id))
	}
	switch cfg_node := cfg_node.(type) {
	case *CFGCondNode:
	}
}
