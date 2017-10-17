package main

import (
	"fmt"
	"math/rand"
	"time"
)

var chGUID chan string

func main() {

	chGUID = make(chan string)
	go GetChannel()
	for i := 0; i < 100; i++ {
		chGUID <- GUID()
	}

}

func GetChannel() {

	for v := range chGUID {
		fmt.Println(v)
	}
}

func GUID() string {
	// generate 32 bits timestamp
	unix32bits := uint32(time.Now().UTC().Unix())

	buff := make([]byte, 12)

	rand.Read(buff)
	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:])
}
