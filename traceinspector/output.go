package traceinspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"traceinspector/imp"
)

type AnalyzerOutputType string

const (
	AnalyzerOutput_error       AnalyzerOutputType = "error"
	AnalyzerOutput_update_node AnalyzerOutputType = "update_node"
	AnalyzerOutput_info        AnalyzerOutputType = "info"
	AnalyzerOutput_warning     AnalyzerOutputType = "warning"
)

type AnalyzerOutput struct {
	Type          AnalyzerOutputType
	Function_name imp.ImpFunctionName
	Node_id       NodeID
	Node_state    string
	Msg           string
}

func write_info(node_location CFGNodeLocation, msg string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(AnalyzerOutput{Type: AnalyzerOutput_info, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}

func write_warning(node_location CFGNodeLocation, msg string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(AnalyzerOutput{Type: AnalyzerOutput_warning, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}

func write_error(node_location CFGNodeLocation, msg string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(AnalyzerOutput{Type: AnalyzerOutput_error, Function_name: node_location.Function_name, Node_id: node_location.Id, Msg: msg})
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
	os.Exit(1)
}

func write_update_node_state(node_location CFGNodeLocation, state_str string, msg string) {
	res := AnalyzerOutput{Type: AnalyzerOutput_update_node, Function_name: node_location.Function_name, Node_id: node_location.Id, Node_state: state_str, Msg: msg}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(res)
	out := &bytes.Buffer{}
	json.Compact(out, buf.Bytes())
	fmt.Println(out.String())
}
