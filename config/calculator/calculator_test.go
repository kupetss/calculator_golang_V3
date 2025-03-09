package calculator

import "fmt"

func test() {
	fmt.Println(Calc("9+1", 0))
	fmt.Println(Calc("(2*5)/2", 0))
	fmt.Println(Calc("1 + (2 + (2 + (2 + 2)))", 0))
	fmt.Println(Calc("2+2*2", 0))
	fmt.Println(Calc("2 / 0", 0))
	fmt.Println(Calc("2.5 * 3", 0))
	fmt.Println(Calc("-1 + 2", 0))
	fmt.Println(Calc("2 * (3 + 4)", 0))
	fmt.Println(Calc("2 * 3 + 4", 0))
	fmt.Println(Calc("2 * (3 + 4 * 2)", 0))
}
