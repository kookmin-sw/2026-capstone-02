//go:build ignore

package main

import "fmt"

func main() {
	x := 0
	fmt.Print("Enter non-negative integer:")
	fmt.Scanf("%d", x)
	for i := x; i > 0; i-- {
		if i%2 == 0 {
			if i%3 == 0 {
				fmt.Print(i, "is a multiple of 2 and 3\n")
			} else {
				fmt.Print(i, "is a multiple of 2 and not a multiple of 3\n")
			}
		} else if i%3 == 0 {
			fmt.Print(i, "is not a multiple of 2, but a multiple of 3 \n")
		} else {
			fmt.Print(i, "is not a multiple of both 2 and 3\n")
		}
	}
	fmt.Print("Done", x)
}
