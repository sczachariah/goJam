package main

import (
	"fmt"
	"goJam/pkg/samples"
)

func main() {
	fmt.Println("Start Main")

	//from hello.go
	samples.HelloWorld()
	//fmt.Println(samples.Swap("Hello", "World"))

	//from duplicate.go
	//samples.StdinDup()
	//samples.FileDup()

	//from fetch.go
	//samples.Fetch()
	//samples.Fetchall()

	//from server.go
	samples.Server()
}
