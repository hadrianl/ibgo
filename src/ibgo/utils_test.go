package ibgo

import (
	"fmt"
	"testing"
)

func TestBytesToInt(t *testing.T) {
	buf := []byte{0, 0, 0, 1}
	size := BytesToInt(buf)
	if size == 1 {
		fmt.Println(size)
	} else {
		t.Errorf("BytesToInt Failed!")
	}

}
