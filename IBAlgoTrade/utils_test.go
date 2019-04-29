package IBAlgoTrade_test

import (
	"fmt"
	"testing"

	. "github.com/hadrianl/ibgo/IBAlgoTrade"
	"github.com/hadrianl/ibgo/ibapi"
)

func TestCreate(t *testing.T) {
	o := new(Trade)
	Create(o)
	fmt.Println(o)
}

func BenchmarkInitDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := new(ibapi.Order)
		ibapi.InitDefault(t)
	}
}
