// con-2 project main.go
package main

import (
	"fmt"
)

func main() {
	go Task1()
	//fmt.Println("Hello World!")
}

func Task1() {
	for i := 0; i < 100; i++ {
		fmt.Println(i)
	}
}
