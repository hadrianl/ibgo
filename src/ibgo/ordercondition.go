package ibgo

import (
	"log"
	"time"
)

type OrderConditioner interface {
	decode(fields [][]byte)
}

type OrderCondition struct {
	CondType                int64
	IsConjunctionConnection bool

	// Price = 1
	// Time = 3
	// Margin = 4
	// Execution = 5
	// Volume = 6
	// PercentChange = 7
}

func (oc OrderCondition) decode(fields [][]byte) {
	connector := decodeString(fields[0])
	oc.IsConjunctionConnection = connector == "a"
}

type ExecutionCondition struct {
	OrderCondition
	SecType  string
	Exchange string
	Symbol   string
}

func (ec ExecutionCondition) decode(fields [][]byte) { // 4 fields
	ec.OrderCondition.decode(fields[0:1])
	ec.SecType = decodeString(fields[1])
	ec.Exchange = decodeString(fields[2])
	ec.Symbol = decodeString(fields[3])
}

type OperatorCondition struct {
	OrderCondition
	IsMore bool
}

func (oc OperatorCondition) decode(fields [][]byte) { // 2 fields
	oc.OrderCondition.decode(fields[0:1])
	oc.IsMore = decodeBool(fields[1])
}

type MarginCondition struct {
	OperatorCondition
	Percent float64
}

func (mc MarginCondition) decode(fields [][]byte) { // 3 fields
	mc.OperatorCondition.decode(fields[0:2])
	mc.Percent = decodeFloat(fields[2])
}

type ContractCondition struct {
	OperatorCondition
	ConId    int64
	Exchange string
}

func (cc ContractCondition) decode(fields [][]byte) { // 4 fields
	cc.OperatorCondition.decode(fields[0:2])
	cc.ConId = decodeInt(fields[2])
	cc.Exchange = decodeString(fields[3])
}

type TimeCondition struct {
	OperatorCondition
	Time time.Time
}

func (tc TimeCondition) decode(fields [][]byte) { // 3 fields
	tc.OperatorCondition.decode(fields[0:2])
	tc.Time = decodeTime(fields[2], "20060102")
}

type PriceCondition struct {
	ContractCondition
	Price         float64
	TriggerMethod int64
}

func (pc PriceCondition) decode(fields [][]byte) { // 6 fields
	pc.ContractCondition.decode(fields[0:4])
	pc.Price = decodeFloat(fields[4])
	pc.TriggerMethod = decodeInt(fields[5])
}

type PercentChangeCondition struct {
	ContractCondition
	ChangePercent float64
}

func (pcc PercentChangeCondition) decode(fields [][]byte) { // 5 fields
	pcc.ContractCondition.decode(fields[0:4])
	pcc.ChangePercent = decodeFloat(fields[4])
}

type VolumeCondition struct {
	ContractCondition
	Volume int64
}

func (vc VolumeCondition) decode(fields [][]byte) { // 5 fields
	vc.ContractCondition.decode(fields[0:4])
	vc.Volume = decodeInt(fields[4])
}

func InitOrderCondition(conType int64) (OrderConditioner, int) {
	var cond OrderConditioner
	var condSize int
	switch conType {
	case 1:
		cond = PriceCondition{}
		condSize = 6
	case 3:
		cond = TimeCondition{}
		condSize = 3
	case 4:
		cond = MarginCondition{}
		condSize = 3
	case 5:
		cond = ExecutionCondition{}
		condSize = 4
	case 6:
		cond = VolumeCondition{}
		condSize = 5
	case 7:
		cond = PercentChangeCondition{}
		condSize = 5
	default:
		log.Panicf("errInitOrderCondition: Unkonwn conType: %v", conType)
	}
	return cond, condSize
}
