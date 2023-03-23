package terra

import "fmt"

func ExampleNumber() {
	n := Number(1)
	fmt.Println(string(n.InternalTokens().Bytes()))
	// 	Output: 1
}
