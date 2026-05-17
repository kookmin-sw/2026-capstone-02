//go:build ignore

package main

func main() {
	a := 1
	b := 2
	c := a + b
	if a+c+b <= 5 {
		c = 10
	} else {
		c = 0
	}

	if b+c <= 0 {
		c = c + 1
	} else {
		c = 10
	}
	a = 1
}
