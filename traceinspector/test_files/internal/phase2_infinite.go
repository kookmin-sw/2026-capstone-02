//go:build ignore

package main

func main() {
	a := 1
	for a > 0 {
		a++
	}
	Print(a, "bob")
	Print(a)
}
