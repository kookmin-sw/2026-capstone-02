package traceinspector

import (
	"encoding/json"
	"os"
	"sort"
	"traceinspector/imp"
)

type AnalyzerOutputType string

const (
	AnalyzerOutput_error       AnalyzerOutputType = "error"
	AnalyzerOutput_update_node AnalyzerOutputType = "update_node"
	AnalyzerOutput_info        AnalyzerOutputType = "info"
	AnalyzerOutput_warning     AnalyzerOutputType = "warning"
)

type AnalyzerOutputHandler struct {
	state_index int `json:"-"`
	Debugs      []AnalyzerOutput
}

type AnalyzerOutput struct {
	Type          AnalyzerOutputType
	State_index   int
	Function_name imp.ImpFunctionName
	Line_number   int
	Node_id       NodeID
	Node_state    string
	Msg           string
}

func (ao *AnalyzerOutputHandler) Print() {
	// buf := &bytes.Buffer{}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	sort.Slice(ao.Debugs, func(i, j int) bool { return ao.Debugs[i].State_index < ao.Debugs[j].State_index })
	enc.Encode(ao)
	// out := &bytes.Buffer{}
	// json.Compact(out, buf.Bytes())
	// fmt.Println(out.String())
}

func (ao *AnalyzerOutputHandler) write_info(node_location CFGNodeLocation, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_info, State_index: ao.state_index,
		Function_name: node_location.Function_name, Line_number: node_location.Line_num, Node_id: node_location.Id, Msg: msg})
}

func (ao *AnalyzerOutputHandler) write_warning(node_location CFGNodeLocation, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_warning, State_index: ao.state_index,
		Function_name: node_location.Function_name, Line_number: node_location.Line_num, Node_id: node_location.Id, Msg: msg})
}

func (ao *AnalyzerOutputHandler) write_error(node_location CFGNodeLocation, msg string) {
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_error, State_index: ao.state_index,
		Function_name: node_location.Function_name, Line_number: node_location.Line_num, Node_id: node_location.Id, Msg: msg})
	ao.Print()
	os.Exit(0)
}

func (ao *AnalyzerOutputHandler) write_update_node_state(node_location CFGNodeLocation, state_str string, msg string) {
	ao.state_index++
	ao.Debugs = append(ao.Debugs, AnalyzerOutput{Type: AnalyzerOutput_update_node, State_index: ao.state_index,
		Function_name: node_location.Function_name, Line_number: node_location.Line_num, Node_id: node_location.Id, Node_state: state_str, Msg: msg})
}
