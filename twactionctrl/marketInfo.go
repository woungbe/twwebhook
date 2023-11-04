package twactionctrl

/*
	공통으로 사용하는 바이낸스 정보를 모와둔 것.
	:


*/

import (
	"fmt"
	binanceorder "twwebhook/binanceOrder"

	"github.com/adshao/go-binance/v2/futures"
)

type MarketInfo struct {
	// 마켓 데이터 symbols
	// 레버리지 데이터 일부분
	Client     binanceorder.BinanceInfo // 오픈데이터를 가져올 임시 데이터 처리
	MarketInfo map[string]futures.Symbol
	Leverage   map[string]*futures.LeverageBracket
}

var MarketInfoPTR *MarketInfo

func GetMarektInfo() *MarketInfo {
	if MarketInfoPTR == nil {
		MarketInfoPTR = new(MarketInfo)
	}
	return MarketInfoPTR
}

// 초기화 작업
func (ty *MarketInfo) Init(acess, key string) {
	ty.Client.Init(acess, key)
	ty.MarketInfo = make(map[string]futures.Symbol)
	ty.Leverage = make(map[string]*futures.LeverageBracket)
}

// 마켓데이터 가져와서 저장하기 - 한번 부르고, 내부데이터 사용하면 됨
func (ty *MarketInfo) GetMarketInfo() error {
	// 마켓데이터 가져오기
	res, err := ty.Client.GetExchangeInfo()
	if err != nil {
		return err
	}

	// 마켓데이터 저장
	for _, v := range res.Symbols {
		ty.MarketInfo[v.Symbol] = v
	}
	return nil
}

// 레버리지 가져오기
func (ty *MarketInfo) GetLeverageBracket() error {
	res, err := ty.Client.GetLeverageList()
	if err != nil {
		return err
	}

	// 마켓데이터 저장
	for _, v := range res {
		ty.Leverage[v.Symbol] = v
	}
	return nil
}

// 심볼의 최대 맥스 값 찾기
func (ty *MarketInfo) GetMaxLeverage(symbol string) (futures.Bracket, error) {
	var tmp futures.Bracket
	leverageList := ty.Leverage[symbol]
	if len(leverageList.Brackets) == 0 {
		return tmp, fmt.Errorf("바스켓에 레버리지 없음")
	}
	return leverageList.Brackets[0], nil
}

// 바스켓 케파 확인하기
func (ty *MarketInfo) GetLeverageBracketDefault(symbol string, leverage int) (futures.Bracket, error) {
	/*

		레버리지 캐파는  50, 20, 10, 5 , 1  이렇게 나열되어 있고, 종목마다 다르다.
		레버리지 25 는 50 케파를 가지고 있다. 21~50까지는 50으로 적용한다.
		레버리지 20 는 11~20 까지 적용된다.

		레버리지 바스켓은 높은 순선대로 나열해서 준다 (바이낸스 기준)
		만약에 다른 거래소라면. 레버리지 값을 높은 순으로 정렬하는 구문을 추가해라.

		현재의 레버리지가 바스켓레버리지 보다 높으면, 이전 것을 준다.  25 라면 -> 50을 준다.
		현재의 레버리지가 바스켓레버리지 와 같다면 현재것을 준다. 20 -> 20 을준다.

	*/
	leverageList := ty.Leverage[symbol]
	var tmp futures.Bracket
	for k, v := range leverageList.Brackets {
		if k == 0 {
			if leverage > v.InitialLeverage {
				return tmp, fmt.Errorf("최대 레버리지보다 높습니다.")
			}
		}
		// 작으면 이전꺼 데이터를 확인
		if leverage > v.InitialLeverage {
			return tmp, nil
		}
		// 같으면 현재 데이터를 확인
		if leverage == v.InitialLeverage {
			return v, nil
		}
		tmp = v
	}

	return tmp, fmt.Errorf("확인된 레버리지가 없습니다. ")
}
