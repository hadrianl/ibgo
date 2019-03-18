package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
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
	format := "20060102 15:04:05"
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

	msgBuf := make([]byte, size)

	if _, err := reader.Read(msgBuf); err != nil {
		return nil, err
	}

	return msgBuf, nil

}

func makeMsgBuf(msg interface{}) []byte {
	switch _msg := msg.(type) {
	case string:
		bs := []byte(_msg)
		return bs
	case int64:
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(_msg))
		return bs
	case float64:
		bs := make([]byte, 8)
		bits := math.Float64bits(_msg)
		binary.LittleEndian.PutUint64(bs, bits)
		return bs
	case []byte:
		return _msg
		// case bool:
		// 	var s string
		// 	if msg {
		// 		s = "1"
		// 	}else {
		// 		s = "0"
		// 	}
		// 	_, err := b.WriteString(s)
		// case time.Time:
		// 	t_string := msg.UTC().Format("20060102 15:04:05"+" UTC")
		// 	_, err := b.WriteString(t_string)
	}
	return nil
}

func mergeMsgBuf(msgBufs [][]byte) []byte {
	msg := bytes.Join(msgBufs, []byte(fieldSplit))
	return msg
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
