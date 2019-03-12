package ibgo

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

// BytesToInt used to convert the first 4 byte into the message size
func BytesToInt(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func panicError(err error) {
	if err != nil {
		panic(err)
		// os.Exit(1)
	}
}

// read_msg try to read the msg based on the message size
func read_msg(reader *bufio.Reader) (int, []byte) {
	sizeBuf, err := reader.Peek(4)
	if err != nil {
		return nil, err
	}
	fmt.Println("error")
	return 1, []byte{1}
}
