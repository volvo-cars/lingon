package terra

import "fmt"

func ExampleBool() {
	b := Bool(true)
	fmt.Println(string(b.InternalTokens().Bytes()))
	// 	Output: true
}
