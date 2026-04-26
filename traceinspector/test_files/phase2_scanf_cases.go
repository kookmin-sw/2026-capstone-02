//go:build ignore

package main

func main() {
	a := 5
	Scanf("%d", a)
	// len_Arr := len(arr)
	if a >= 0 {
		a = 1
	} else {
		a = -1
	}
	Print(a)
}
