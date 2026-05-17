//go:build ignore

package main

func main() {
	a := 1
	if a > 0 { //true
		if a <= 0 { // false
			Print("unreachable")
		}
	}

	b := 5

	if a < b { // true
		a++
	}
	if a < b { // true
		Print("always true")
	}

	for a < b {
		a++
	}
	if a < b {
		Print("unreachable")
	}

	for a > 0 {
		a++
	}
	Print(a)
}
