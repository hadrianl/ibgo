package main

import (
	"fmt"
	"ibgo"
)

func main() {
	client := ibgo.NewIbClient("127.0.0.1", 7497, 6)

	buf := []byte{0, 1, 1, 1}

	fmt.Println(ibgo.BytesToInt(buf))

}
