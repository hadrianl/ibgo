package main

import (
	"bytes"
	"fmt"
)

func main() {
	f := ""
	b := []byte(f)
	c := []byte(" ")
	fmt.Println(len(bytes.Split(b, c)))
	fmt.Println("string:", f, "byte:", b, c)
}
