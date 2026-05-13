//go:build ignore

package main

import "fmt"

func main() {
	x := 0
	fmt.Print("Enter a number: ")
	fmt.Scanf("%d", x)
	if x < 0 {
		x = -x + 100
		fmt.Print("got negative\n")
	} else {
		x = -x
	}
	fmt.Print("result is", x, "\n")
}
