//go:build ignore

package main

import "fmt"

func subone(a int) int {
	return a - 1
}

func factorial_naive(a int) int {
	if a <= 1 {
		return a
	}
	return a * factorial_naive(subone(a))
}

func main() {
	x := 0
	fmt.Print("Enter a number: ")
	fmt.Scanf("%d", x)
	fmt.Print("Factorial of", x, "is", factorial_naive(x))
}
