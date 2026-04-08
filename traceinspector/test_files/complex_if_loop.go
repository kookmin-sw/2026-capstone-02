//go:build ignore

package main

import "fmt"

func insertionSort(arr []int) []int {
	i := 1

	for i < len(arr) {
		key := arr[i]
		j := i - 1

		for j >= 0 {
			if arr[j] > key {
				arr[j+1] = arr[j]
				j = j - 1
				continue
			}

			break
		}

		arr[j+1] = key
		i = i + 1
	}

	return arr
}

func chaos(arr []int) int {
	i := 0
	result := 0

	for i < len(arr) {
		if arr[i] == 0 {
			i = i + 1
			continue
		}

		j := 0
		for j < arr[i] {
			k := 0

			for k < j {
				if k%2 == 0 {
					k = k + 1
					continue
				}

				if k > 3 {
					break
				}

				result = result + k
				k = k + 1
			}

			if j > 5 {
				break
			}

			j = j + 1
		}

		i = i + 1
	}

	return result
}

func findFirstMatch(matrix [][]int, target int) int {
	i := 0

	for i < len(matrix) {
		j := 0

		for j < len(matrix[i]) {
			if matrix[i][j] == target {
				return i*100 + j
			}

			if matrix[i][j] < 0 {
				break
			}

			j = j + 1
		}

		i = i + 1
	}

	return -1
}

func countSpecial(arr []int) int {
	count := 0
	i := 0

	for i < len(arr) {
		j := 0

		for j < arr[i] {
			if j%2 == 0 {
				j = j + 1
				continue
			}

			if j > 5 {
				break
			}

			count = count + 1
			j = j + 1
		}

		if arr[i] == 0 {
			break
		}

		i = i + 1
	}

	return count
}

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
	fmt.Print("Done", x, "\n")
	arr := []int{8, 7, 6, 5, 4, 3, 2, 1}
	fmt.Print(insertionSort(arr))
}
