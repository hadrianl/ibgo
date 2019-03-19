package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	fieldSplit string = "\x00"
)

// bytesToInt used to convert the first 4 byte into the message size
func bytesToInt(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func bytesToTime(buf []byte) time.Time {
	format := "20060102 15:04:05 CST"
	t := string(buf)
	localtime, err := time.ParseInLocation(format, t, time.Local)
	if err != nil {
		fmt.Println(err)
	}
	return localtime
}

func panicError(err error) {
	if err != nil {
		panic(err)
		// os.Exit(1)
	}
}

// readMsg try to read the msg based on the message size
func readMsgBuf(reader *bufio.Reader) ([]byte, error) {
	sizeBuf := make([]byte, 4)

	if _, err := reader.Read(sizeBuf); err != nil {
		return nil, err
	}
	size := bytesToInt(sizeBuf)
	// fmt.Println("Get SizeBuf:", size)

	msgBuf := make([]byte, size)

	if _, err := reader.Read(msgBuf); err != nil {
		return nil, err
	}

	return msgBuf, nil

}

func makeMsgBuf(msg interface{}) []byte {
	return append([]byte(fmt.Sprint(msg)), '\x00')
	// switch msg.(type) {
	// case string:
	// 	bs := []byte(msg.(string))
	// 	return append(bs, '\x00')

	// case int64:
	// 	bs := make([]byte, 8)
	// 	binary.BigEndian.PutUint64(bs, uint64(msg.(int64)))
	// 	return append(bs, '\x00')
	// case float64:
	// 	bs := make([]byte, 8)
	// 	bits := math.Float64bits(msg.(float64))
	// 	binary.LittleEndian.PutUint64(bs, bits)
	// 	return append(bs, '\x00')
	// case []byte:
	// 	return append(msg.([]byte), '\x00')
	// 	// case bool:
	// 	// 	var s string
	// 	// 	if msg {
	// 	// 		s = "1"
	// 	// 	}else {
	// 	// 		s = "0"
	// 	// 	}
	// 	// 	_, err := b.WriteString(s)
	// 	// case time.Time:
	// 	// 	t_string := msg.UTC().Format("20060102 15:04:05"+" UTC")
	// 	// 	_, err := b.WriteString(t_string)
	// }
	// return nil
}

func mergeMsgBuf(fields ...interface{}) []byte {
	msgBufs := [][]byte{}
	for _, f := range fields {
		msgBufField := makeMsgBuf(f)
		msgBufs = append(msgBufs, msgBufField)
	}
	msg := bytes.Join(msgBufs, []byte(""))
	return msg
}

func makeMsg(fields ...interface{}) []byte {
	msg := mergeMsgBuf(fields...)
	sizeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuf, uint32(len(msg)))

	return append(sizeBuf, msg...)
}

func splitMsgBuf(data []byte) [][]byte {
	fields := bytes.Split(data, []byte(fieldSplit))
	return fields[:len(fields)-1]

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
