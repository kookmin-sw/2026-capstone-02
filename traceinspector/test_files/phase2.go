//go:build ignore

package main

func main() {
	a := 1
	for i := 0; i <= 10; i++ {
		Print(i)
		a = i
	}
	a = i
	Print(a, "bob")
	Print(a)
}
