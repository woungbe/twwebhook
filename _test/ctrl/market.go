package ctrl

import (
	binanceorder "twwebhook/binanceOrder"

	"github.com/adshao/go-binance/v2/futures"
)

/*
 마켓 정보를 다루는 루틴
 바이낸스 테스트 정보를 가져와야지
*/

type BianaceExchangeInfo struct {
	C            binanceorder.BinanceInfo
	ExInfo       []*futures.Symbol
	LeverageList []*futures.LeverageBracket
}

// Init 처리
func (ty BianaceExchangeInfo) Init(key, secrit string) {
	var client binanceorder.BinanceInfo
	ty.C = client
	ty.C.Init(key, secrit)
}

// 마켓정보
func (ty BianaceExchangeInfo) GetMarketInfo() error {
	// (*futures.ExchangeInfo, error)
	exchangeinf, err := ty.C.GetExchangeInfo()
	if err != nil {
		return err
	}
	ty.ExInfo = exchangeinf.Symbols
}

// 레버리지 바스켓 정보
func (ty BianaceExchangeInfo) GetLeverageList() error {
	LeverageBracket, err := ty.C.GetLeverageList()
	if err != nil {
		return err
	}
	ty.LeverageList = LeverageBracket
}

// 해당 심볼의 가격에 대한 소수점 자릿수 리턴해주기
func (ty BianaceExchangeInfo) GetPriceTickSize(symbol string) int {
	// 심볼을 찾아서 틱사이즈 리턴
	for k, v := range ty.ExInfo {
		if v.Symbol == symbol {
			// 탁사이즈 정리
			return 5
		}
	}
}

// 해당심볼의 가격에 대한 소수점 정해서 전달해주기
func (ty BianaceExchangeInfo) CheckPriceTickSize(symbol, price string) (string, error) {
	ticksize := GetPriceTickSize(symbol)

	// 현재 호출된 가격과, ticksize 정리해서 리턴해주기
	return "", nil
}

// 수량 체크
func (ty BianaceExchangeInfo) GetCoinTickSize(symbol string) (string, error) {
	// 최소 수량 사이즈에 대해서 리턴
	for k, v := range ty.ExInfo {
		if v.Symbol == symbol {
			// 탁사이즈 정리
			return 5
		}
	}
}

// 수량 소수점 체크 리턴
func (ty BianaceExchangeInfo) CheckCoinTickSize(symbol, price string) (string, error) {
	lotsize := ty.GetCoinTickSize(symbol)

	//
	return "", nil
}

// 금액 - 최소/최대 체크
func (ty BianaceExchangeInfo) GetMinMaxAmount(symbol string) (string, string) {
	// lotsize := ty.GetCoinTickSize(symbol)

	// 최소 금액 리턴
	// 최대 금액 리턴
	for k, v := range ty.ExInfo {
		if v.Symbol == symbol {
			// 탁사이즈 정리
			return "5", "100000"
		}
	}

	return 0, 0
}

// 수량 - 최소/최대 체크
func (ty BianaceExchangeInfo) GetMinMaxQty(symbol string) (string, string) {
	// lotsize := ty.GetCoinTickSize(symbol)

	// 최소 금액 리턴
	// 최대 금액 리턴
	for k, v := range ty.ExInfo {
		if v.Symbol == symbol {
			// 탁사이즈 정리
			return "5", "100000"
		}
	}

	return 0, 0
}

// 수량도 만족하고, 금액도 만족하는 최소 수량 체크하기
func (ty BianaceExchangeInfo) GetMinAmount(symbol string, price string) string {

	// 주문가능한 최소 수량 찾기

	// 가격도 최소 수량 이상일때
	// 수량이 최소 수량 이상이면서

	// 일단 가격을 넘는 최소 수량 1.02

	// 해당 수량이 최소수량 을 넘으면 오케이

	// 해당 수량이 최소수량이 안된다면, 최소 수량 기준으로 재정리

	// 최소 총금액, 최대 총금액
	minPrice, maxPrice := ty.GetMinMaxAmount(symbol)

	// 최소 수량, 최대 수량
	minQty, maQty := ty.GetMinMaxQty(symbol)

	// 주문 최소 수량
	qty := (minPrice / price) + minQty

	if qty <= minQty {
		qty = minQty
	}

	// 최소금액을 만족하는 최소 수량
	return qty
}
