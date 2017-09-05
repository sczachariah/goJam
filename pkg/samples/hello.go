package samples

import (
	"fmt"
)

func HelloWorld() {
	fmt.Println("Hello World")
}

func Swap(a, b string) (string, string) {
	return b, a
}
