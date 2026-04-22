package traceinspector

import (
	"fmt"
	"traceinspector/domain"
	"traceinspector/imp"
)

// This is the main struct for initializing abstract interpretation.
//
// abstract_semantics: The interface implementing the abstract step relation
// function_mem_map:
type AbstractAnalyzer[IntDom domain.IntegerDomain[IntDom], ArrDom ArrayDomain[IntDom, ArrDom]] struct {
	abstract_semantics AbstractSemantics[IntDom, ArrDom]
	function_mem_map   map[imp.ImpFunctionName]*FunctionAbstractMem[IntDom, ArrDom]
	function_cfgs      FunctionCFGMap
	function_defs      imp.ImpFunctionMap
}

func (analyzer *AbstractAnalyzer[IntDomainImpl, ArrayDomainImpl]) Start_analysis(function_name imp.ImpFunctionName) {
	analyzer.function_mem_map = make(map[imp.ImpFunctionName]*FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl])
	analyzer.function_mem_map[function_name] = &FunctionAbstractMem[IntDomainImpl, ArrayDomainImpl]{}
	analyzer.function_mem_map[function_name].Initialize(function_name)

	initial_state := AbstractState[IntDomainImpl, ArrayDomainImpl]{node_location: analyzer.function_cfgs[function_name].Entry_node, abstract_mem: make(AbstractNodeMem[IntDomainImpl, ArrayDomainImpl])}
	for _, val := range analyzer.abstract_semantics.Step(initial_state) {
		fmt.Println(val)
	}
}

func Test(func_cfg_map FunctionCFGMap, func_name imp.ImpFunctionName, func_info_map imp.ImpFunctionMap) {
	func_mem := FunctionAbstractMem[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{
		mem:           make(map[NodeID]AbstractNodeMem[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]),
		function_name: func_name,
		return_value:  AbstractValue[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{},
	}
	semantics := ImpFunctionInterpreter[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{
		func_cfg_map:        func_cfg_map,
		func_name:           func_name,
		func_info_map:       func_info_map,
		abstract_mem:        &func_mem,
		intdomain_default:   domain.IntervalDomain{},
		booldomain_default:  domain.BoolDomain{},
		arraydomain_default: ArraySummaryDomain[domain.IntervalDomain]{},
	}
	g := AbstractAnalyzer[domain.IntervalDomain, ArraySummaryDomain[domain.IntervalDomain]]{abstract_semantics: &semantics, function_cfgs: func_cfg_map, function_defs: func_info_map}
	g.Start_analysis("main")
}
