package main

import (
	"fmt"
	"ibgo"
	"time"
)

func main() {
	var err error
	ic := &ibgo.IbClient{}
	err = ic.Connect("127.0.0.1", 7497, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("afterConnect")
	err = ic.HandShake()
	if err != nil {
		panic(err)
	}
	fmt.Println("afterhandShake")

	ic.Run()
	time.Sleep(10 * time.Second)

}
