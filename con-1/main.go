// con-1 project main.go
package main

import (
	"fmt"
	"runtime"
	_ "time"
)

func main() {
	runtime.GOMAXPROCS(4)
	for i := 0; i <= 100; i++ {
		go fmt.Println("Hello World!->", i)
		//go fmt.Println(time.Now())
		go func() {
			//go fmt.Println(time.Now())
			fmt.Println("Hello Another world--->", i)
		}()
		runtime.Gosched()
	}
}
