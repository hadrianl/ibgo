package main

// func main() {
// 	var err error
// 	ic := &ibgo.IbClient{}
// 	err = ic.Connect("127.0.0.1", 7497, 10)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("afterConnect")
// 	err = ic.HandShake()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("afterhandShake")

// 	ic.Run()
// 	time.Sleep(10 * time.Second)

// }

// func main() {
// 	// var msg bytes.Buffer
// 	// head := []byte("API\x00")
// 	// // minVer := []byte("100")
// 	// // maxVer := []byte("148")
// 	// // connectOptions := []byte("")
// 	// clientVersion := []byte("v100..148")
// 	// sizeofCV := make([]byte, 4)
// 	// binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))

// 	// msg.Write(head)
// 	// msg.Write(sizeofCV)
// 	// msg.Write(clientVersion)
// 	// fmt.Println(msg.Bytes())
// 	// f := ibgo.Split_msg([]byte("API\x00sfsdfs\x00dfsfs\x00"))
// 	// fmt.Println(f)
// 	var clientId int = 10
// 	v := make([]byte, 8)
// 	binary.BigEndian.PutUint64(v, uint64(clientId))
// 	fmt.Println(v)
// 	fmt.Println(bytes.Join([][]byte{v, []byte("\x00")}, []byte("")))
// }
