package terra

import "fmt"

func ExampleString() {
	s := String("hello world")
	fmt.Println(string(s.InternalTokens().Bytes()))
	// 	Output: "hello world"
}
