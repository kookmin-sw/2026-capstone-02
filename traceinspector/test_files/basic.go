//go:build ignore

package main

func add1(a int) int {
	return a + 1
}

func main() {
	a := add1(1)
	_ = a
	b, c := a, 2
	_ = b + c
}
