package binanceorder

import (
	"fmt"
	"testing"
)

var access string = "MhPq509wvXhgdx2EXJ3gO132MRMMByCCHNrCW9p6apgYvlgzueyzmbByIgfO34ah"
var secrit string = "dp3k5cllVlayPH3Eb5HrTSA4ktDovurFM2rWPkU6qrSHDI9HpcfpqU5LJyM8NWf8"

// 되는것만 확인했으면 됐음
func TestGetExchangeInfo(t *testing.T) {
	var S BinanceInfo
	S.Init(access, secrit)

	res, err := S.GetExchangeInfo()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestGetLeverageList(t *testing.T) {
	var S BinanceInfo
	S.Init(access, secrit)

	res, err := S.GetLeverageList()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		fmt.Printf("%+v\n", v)
	}
}

func TestGetUserLevelage(t *testing.T) {
	var S BinanceInfo
	S.Init(access, secrit)

	res, err := S.GetPositionInfo("BTCUSDT")
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range res {
		fmt.Printf("%+v\n", v)
	}

}
