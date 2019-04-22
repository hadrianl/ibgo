package ibgo

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

const (
	fieldSplit  byte    = '\x00'
	UNSETFLOAT  float64 = math.MaxFloat64
	UNSETINT    int64   = math.MaxInt64
	NO_VALID_ID int64   = -1
)

// // bytesToInt used to convert the first 4 byte into the message size
// func bytesToInt(buf []byte) int32 {
// 	return int32(binary.BigEndian.Uint32(buf))
// }

func bytesToTime(buf []byte) time.Time {
	format := "20060102 15:04:05 CST"
	t := string(buf)
	localtime, err := time.ParseInLocation(format, t, time.Local)
	if err != nil {
		log.Println(err)
	}
	return localtime
}

// readMsg try to read the msg based on the message size
func readMsgBuf(reader *bufio.Reader) ([]byte, error) {
	sizeBuf := make([]byte, 4)

	// if n, err := reader.Read(sizeBuf); err != nil || n != 4 {
	// 	panic(fmt.Errorf("errReadSize: read %v err: %v", n, err))
	// 	return nil, err
	// }

	for n, r := 0, 4; n < r; {
		tempMsgBuf := make([]byte, r-n)
		tn, err := reader.Read(tempMsgBuf)
		if err != nil {
			return nil, err
		}

		copy(sizeBuf[n:n+tn], tempMsgBuf)
		n += tn
	}
	// log.Println("Get SizeBuf:", sizeBuf)
	size := int(binary.BigEndian.Uint32(sizeBuf))
	// log.Println("Get SizeBuf:", size)

	msgBuf := make([]byte, size)

	// XXX: maybe there is a better way to get fixed size of bytes
	for n, r := 0, size; n < r; {
		tempMsgBuf := make([]byte, r-n)
		tn, err := reader.Read(tempMsgBuf)
		if err != nil {
			return nil, err
		}

		copy(msgBuf[n:n+tn], tempMsgBuf)
		n += tn

	}

	return msgBuf, nil

}

func field2Buf(msg interface{}) []byte {
	var b []byte

	switch msg.(type) {

	case int:
		b = encodeInt(msg.(int))
	case int64:
		b = encodeInt64(msg.(int64))
	case OUT:
		b = encodeInt64(int64(msg.(OUT))) // maybe there is a better solution
	case float64:
		b = encodeFloat(msg.(float64))
	case string:
		b = encodeString(msg.(string))
	case bool:
		b = encodeBool(msg.(bool))
	case []byte:
		b = msg.([]byte)
	case time.Time:
		b = encodeTime(msg.(time.Time))

	default:
		log.Panic("errmakeMsgBuf: can't converst the param")
	}

	return append(b, '\x00')
}

// makeMsgBuf is a universal way to make the request ,but not an efficient way
// TODO: do some test and improve!!!
func makeMsgBuf(fields ...interface{}) []byte {

	// make the whole msg buf
	msgBufs := [][]byte{}
	for _, f := range fields {
		// make the field into buf
		msgBuf := field2Buf(f)
		msgBufs = append(msgBufs, msgBuf)
	}
	msg := bytes.Join(msgBufs, []byte(""))

	// add the size header
	sizeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuf, uint32(len(msg)))

	return append(sizeBuf, msg...)
}

func splitMsgBuf(data []byte) [][]byte {
	fields := bytes.Split(data, []byte{fieldSplit})
	return fields[:len(fields)-1]

}

func decodeInt(field []byte) int64 {
	if bytes.Equal(field, []byte{}) {
		return 0
	}
	i, err := strconv.ParseInt(string(field), 10, 64)
	if err != nil {
		log.Panicf("errDecodeInt: %v", err)
	}
	return i
}

func decodeIntCheckUnset(field []byte) int64 {
	if bytes.Equal(field, []byte{}) {
		return math.MaxInt64
	}
	i, err := strconv.ParseInt(string(field), 10, 64)
	if err != nil {
		log.Panicf("errDecodeInt: %v", err)
	}
	return i
}

func decodeFloat(field []byte) float64 {
	if bytes.Equal(field, []byte{}) || bytes.Equal(field, []byte("None")) {
		return 0.0
	}

	f, err := strconv.ParseFloat(string(field), 64)
	if err != nil {
		log.Panicf("errDecodeFloat: %v", err)
	}

	return f
}

func decodeFloatCheckUnset(field []byte) float64 {
	if bytes.Equal(field, []byte{}) || bytes.Equal(field, []byte("None")) {
		return math.MaxFloat64
	}

	f, err := strconv.ParseFloat(string(field), 64)
	if err != nil {
		log.Panicf("errDecodeFloat: %v", err)
	}

	return f
}

func decodeBool(field []byte) bool {

	if bytes.Equal(field, []byte{'0'}) || bytes.Equal(field, []byte{}) {
		return false
	}
	return true
}

func decodeString(field []byte) string {
	return string(field)
}

func decodeDate(field []byte) time.Time {
	if len(field) != 8 || bytes.Equal(field, []byte{}) {
		return time.Time{}
	}
	tstring := string(field)
	t, err := time.Parse("20060102", tstring)
	if err != nil {
		log.Printf("errDeocodeDate: %v  tstring: %v", err, tstring)
		return time.Time{}
	}

	return t
}

func decodeTime(field []byte, layout string) time.Time {
	if bytes.Equal(field, []byte{}) {
		return time.Time{}
	}

	t, err := time.Parse(layout, string(field))
	if err != nil {
		log.Panicf("errDeocodeTime: %v  format: %v", field, layout)
	}
	return t
}

func encodeInt64(i int64) []byte {
	bs := []byte(strconv.FormatInt(i, 10))
	return bs
}

func encodeInt(i int) []byte {
	bs := []byte(strconv.Itoa(i))
	return bs
}

func encodeFloat(f float64) []byte {
	bs := []byte(strconv.FormatFloat(f, 'g', 10, 64))
	return bs
}

func encodeString(str string) []byte {
	bs := []byte(str)
	return bs
}

func encodeBool(b bool) []byte {
	if b {
		return []byte{'1'}
	}
	return []byte{'0'}

}

func encodeTagValue(tv TagValue) []byte {
	return []byte(fmt.Sprintf("%v=%v;", tv.Tag, tv.Value))
}

func encodeTime(t time.Time) []byte {
	return []byte{} // TODO
}

func handleEmpty(d interface{}) string {
	switch d.(type) {
	case int64:
		v := d.(int64)
		if v == UNSETINT {
			return ""
		} else {
			return strconv.FormatInt(v, 10)
		}
	case float64:
		v := d.(float64)
		if v == UNSETFLOAT {
			return ""
		} else {
			return strconv.FormatFloat(v, 'g', 10, 64)
		}
	default:
		panic("handleEmpty error")

	}
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
