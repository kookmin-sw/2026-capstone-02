//go:build ignore

package main

import "fmt"

func add1(a int) int {
	return a + 1
}

func composite(a int, b int) bool {
	g := (a-b)-(0-1) <= 2                    // a + -b <= 1
	return 5*a+4*2-2+(2+a)*b+11+add1(a) == 1 // 5a + 8 - 2 + 2b + ab + 11 + a + 1 = 6a + ab + 2b + 18
}

func main() {
	a := add1(1)
	b := add1(a)
	c := composite(a, b)
	fmt.Print(b)
}
