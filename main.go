package main

import (
	"fmt"
)

type testingStruct struct {
	A  int64
	BS *B
}
type B struct {
	C int64
}

func main() {
	// f := ""
	// b := []byte(f)
	// c := []byte(" ")
	// fmt.Println(len(bytes.Split(b, c)))
	// fmt.Println("string:", f, "byte:", b, c)
	var t []testingStruct
	fmt.Println(t)
	fmt.Println(len(t))
	fmt.Println(t == nil)
	t = []testingStruct{}
	fmt.Println(t == nil)
	// fmt.Println(new(testingStruct) == nil)
	// fmt.Println(&testingStruct{} == nil)
	// var ts *testingStruct
	// fmt.Println(ts == nil)
	// var t *testingStruct
	// fmt.Println(t == nil)
	// fmt.Println(t.BS)
	// fmt.Println(t.BS == nil)
}
