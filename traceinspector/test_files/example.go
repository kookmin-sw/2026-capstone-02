package testfiles

import "fmt"

func main() {
	var x int
	fmt.Scan(&x)
	if x < 0 {

	} else {
		x = -x
	}
	fmt.Println(x)
}
