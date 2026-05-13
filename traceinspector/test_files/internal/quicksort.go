//go:build ignore

package main

import "fmt"

func partition(a []int, lo int, hi int) int {
	pivot := a[hi]
	i := lo
	j := lo

	for j < hi {
		if a[j] < pivot {
			// swap a[i], a[j]
			tmp := a[i]
			a[i] = a[j]
			a[j] = tmp

			i = i + 1
		}
		j = j + 1
	}

	// place pivot
	tmp := a[i]
	a[i] = a[hi]
	a[hi] = tmp

	return i
}

func quicksort(a []int, lo int, hi int) []int {
	if lo >= hi {
		return a
	}

	p := partition(a, lo, hi)

	quicksort(a, lo, p-1)
	quicksort(a, p+1, hi)
	return a
}

func main() {
	n := 0
	fmt.Scanf("%d", n)
	arr := make_array(n, 0)
	for i := 0; i < n; i++ {
		x := 0
		fmt.Scanf("%d", x)
		arr[i] = x
	}
	fmt.Print(quicksort(arr, 0, n-1), "\n")
}
