//go:build ignore

package main

func main() {
	a := 1
	if a > 0 {
		if a <= 0 {
			Print("unreachable")
		}
	}

	for a > 0 {
		a++
	}
	Print(a)
}
