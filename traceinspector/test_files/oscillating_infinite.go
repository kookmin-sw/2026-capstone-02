//go:build ignore

package main

func main() {
	a := 1
	b := 1
	for true {
		a = a + b
		a = -a
		b = -b
		Print(a) // a = 2, -3, 4, -5, 6, -7, ...
	}
	Print(a)
}
