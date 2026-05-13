//go:build ignore

package main

func main() {
	a := 1
	for true {
		a++
		if a > 5 {
			break
		}
	}
	Print(a, "bob")
	Print(a)
}
