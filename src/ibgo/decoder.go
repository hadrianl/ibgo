package ibgo

import (
	"strconv"
	"strings"
)

type IbDecoder struct {
	wrapper       *IbWrapper
	version       Version
	msgId2process map[IN]func([][]byte)
}

//NewIbDecoder create a decoder to decode the fileds based on version
func NewIbDecoder(wrapper *IbWrapper, version Version) *IbDecoder {
	decoder := IbDecoder{}
	decoder.wrapper = wrapper
	decoder.version = version
	return &decoder
}

func (d *IbDecoder) setVersion(version Version) {
	d.version = version
}

func (d *IbDecoder) interpret(fs ...[]byte) {
	if len(fs) == 0 {
		return
	}

	MsgId, _ := strconv.ParseInt(string(fs[0]), 10, 64)
	processer := d.msgId2process[IN(MsgId)]
	processer(fs[1:])
}

func (d *IbDecoder) setmsgId2process() {
	d.msgId2process = map[IN]func([][]byte){

		NEXT_VALID_ID: d.processNextValidId,
		MANAGED_ACCTS: d.processManagedAccounts,
	}

}

func (d *IbDecoder) processNextValidId(f [][]byte) {
	reqId, _ := strconv.Atoi(string(f[1]))
	d.wrapper.nextValidId(reqId)

}

func (d *IbDecoder) processManagedAccounts(f [][]byte) {
	accNames := strings.Split(string(f[1]), ",")

	accsList := []Account{}
	for _, acc := range accNames {
		accsList = append(accsList, Account{Name: acc})
	}
	d.wrapper.managedAccounts(accsList)

}

func (d *IbDecoder) processTickPriceMsg(f [][]byte) {

}

func (d *IbDecoder) processOrderStatusMsg(f [][]byte) {

}

func (d *IbDecoder) processOpenOrder(f [][]byte) {

}
func (d *IbDecoder) processPortfolioValueMsg(f [][]byte) {

}
func (d *IbDecoder) processContractDataMsg(f [][]byte) {

}
func (d *IbDecoder) processBondContractDataMsg(f [][]byte) {

}
func (d *IbDecoder) processScannerDataMsg(f [][]byte) {

}
func (d *IbDecoder) processExecutionDataMsg(f [][]byte) {

}
func (d *IbDecoder) processHistoricalDataMsg(f [][]byte) {

}
func (d *IbDecoder) processHistoricalDataUpdateMsg(f [][]byte) {

}
func (d *IbDecoder) processRealTimeBarMsg(f [][]byte) {

}
func (d *IbDecoder) processTickOptionComputationMsg(f [][]byte) {

}

func (d *IbDecoder) processDeltaNeutralValidationMsg(f [][]byte) {

}
func (d *IbDecoder) processMarketDataTypeMsg(f [][]byte) {

}
func (d *IbDecoder) processCommissionReportMsg(f [][]byte) {

}
func (d *IbDecoder) processPositionDataMsg(f [][]byte) {

}
func (d *IbDecoder) processPositionMultiMsg(f [][]byte) {

}
func (d *IbDecoder) processSecurityDefinitionOptionParameterMsg(f [][]byte) {

}
func (d *IbDecoder) processSecurityDefinitionOptionParameterEndMsg(f [][]byte) {

}
func (d *IbDecoder) processSoftDollarTiersMsg(f [][]byte) {

}
func (d *IbDecoder) processFamilyCodesMsg(f [][]byte) {

}
func (d *IbDecoder) processSymbolSamplesMsg(f [][]byte) {

}
func (d *IbDecoder) processSmartComponents(f [][]byte) {

}
func (d *IbDecoder) processTickReqParams(f [][]byte) {

}
func (d *IbDecoder) processMktDepthExchanges(f [][]byte) {

}

func (d *IbDecoder) processHeadTimestamp(f [][]byte) {

}
func (d *IbDecoder) processTickNews(f [][]byte) {

}
func (d *IbDecoder) processNewsProviders(f [][]byte) {

}
func (d *IbDecoder) processNewsArticle(f [][]byte) {

}
func (d *IbDecoder) processHistoricalNews(f [][]byte) {

}
func (d *IbDecoder) processHistoricalNewsEnd(f [][]byte) {

}
func (d *IbDecoder) processHistogramData(f [][]byte) {

}
func (d *IbDecoder) processRerouteMktDataReq(f [][]byte) {

}
func (d *IbDecoder) processRerouteMktDepthReq(f [][]byte) {

}
func (d *IbDecoder) processMarketRuleMsg(f [][]byte) {

}
func (d *IbDecoder) processPnLMsg(f [][]byte) {

}
func (d *IbDecoder) processPnLSingleMsg(f [][]byte) {

}
func (d *IbDecoder) processHistoricalTicks(f [][]byte) {

}
func (d *IbDecoder) processHistoricalTicksBidAsk(f [][]byte) {

}
func (d *IbDecoder) processHistoricalTicksLast(f [][]byte) {

}
func (d *IbDecoder) processTickByTickMsg(f [][]byte) {

}
func (d *IbDecoder) processOrderBoundMsg(f [][]byte) {

}
func (d *IbDecoder) processMarketDepthL2Msg(f [][]byte) {

}

// ----------------------------------------------------
