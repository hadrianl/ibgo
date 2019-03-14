package ibgo

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type version int

const (
	MIN_CLIENT_VER version = 100
	MAX_CLIENT_VER version = 148
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
	buf := []byte{}
	_, err := reader.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("error")
	return 1, []byte{1}
}

// func ibWrite(b *bytes.Buffer, msg interface{}) error {
// 	switch reflect.TypeOf(msg) {
// 	case string:
// 		_, err := b.WriteString(msg.(string) + "\000")
// 	case int64:
// 		_, err := b.WriteString(strconv.FormatInt(msg.(int64), 10))
// 	case float64:
// 		_, err := b.WriteString(strconv.FormatFloat(msg.(float64), 'g', 10, 63))
// 	case []byte:
// 		_, err := b.Write(msg)
// 	case bool:
// 		var s string
// 		if msg {
// 			s = "1"
// 		}else {
// 			s = "0"
// 		}
// 		_, err := b.WriteString(s)
// 	case time.Time:
// 		t_string := msg.UTC().Format("20060102 15:04:05"+" UTC")
// 		_, err := b.WriteString(t_string)
// 	}

// 	return err
// }
