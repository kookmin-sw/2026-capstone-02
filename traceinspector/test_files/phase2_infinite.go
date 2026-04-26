//go:build ignore

package main

func main() {
	a := 1
	for a <= 50 {
		a++
	}
	Print(a, "bob")
	Print(a)
}
