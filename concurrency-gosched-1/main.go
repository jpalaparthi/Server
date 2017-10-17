// concurrency-gosched-1 project main.go
package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func main() {
	{
		fmt.Println("CPU number", runtime.NumCPU())
		//runtime.GC()
		runtime.GOMAXPROCS(0)
		iterations := 10
		for i := 0; i < iterations; i++ {
			go display(i)
		}
	}

	fmt.Println("Good bye")
	_, _ = fmt.Scanln()
	runtime.Gosched()
}
func display(num int) {
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println(num, ":", t)
	time.Sleep(time.Microsecond * 10)
}
