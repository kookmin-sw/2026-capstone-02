//go:build ignore

package main

import "fmt"

func main() {
	x := 0
	for true {
		x++
		if x > 5 {
			break
		}
	}
	fmt.Print("Done", x)
}
