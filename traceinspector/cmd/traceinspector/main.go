package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"traceinspector"
)

func main() {
	// https://pkg.go.dev/flag#String
	input_path := flag.String("gofile", "", "")
	_ = flag.Bool("print-cfg", false, "whether to just print cfg and false")
	flag.Parse()
	if *input_path == "" {
		panic("need to pass input go file path with --gofile")
	}

	just_print_cfg := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "print-cfg" {
			just_print_cfg = true
		}
	})

	fset := token.NewFileSet()
	// https://pkg.go.dev/go/parser#Mode
	file, err := parser.ParseFile(fset, *input_path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error while parsing", input_path, "-", err)
		return
	}

	if just_print_cfg {
		traceinspector.Print_cfg(file, fset)
		return
	}
}
