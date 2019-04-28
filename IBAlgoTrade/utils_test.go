package IBAlgoTrade_test

import (
	"fmt"
	"testing"

	. "github.com/hadrianl/ibgo/IBAlgoTrade"
	"github.com/hadrianl/ibgo/ibapi"
)

func TestCreate(t *testing.T) {
	o := new(ibapi.Order)
	Create(o)
	fmt.Println(o)
}
