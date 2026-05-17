//go:build ignore

package main

func main() {
	a := 1
	for true {
		a++
		if a > 5 { // a must be 6
			break
		}
	} // a must be 6
	Print(a)
}
