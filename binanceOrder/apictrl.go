package binanceorder

import (
	"fmt"
	"strconv"
	"twwebhook/utils"

	"github.com/adshao/go-binance/v2/futures"
)

/*
	라이브러리 주의사항 - 한번 한 작업을 다시하지 않게 하는 작업
	: 라이브러리는 쌓아 가는 것.  그래야 많은 일을 하고도 시간이 남음

	주문할때마다 체크하는거 귀찮아서 만듬.
	처음에만 주문만 만들고,
	두번째는 소켓도 연결해야지.

*/

type BiananceOrderCtrl struct {
	MarketList map[string]futures.Symbol          // 공통으로 사용할 마켓데이터
	Leverage   map[string]futures.LeverageBracket // 공통으로 사용할 소켓 데이터

	B BinanceInfo // 바이낸스 api 콜
}

var BianancePTR *BiananceOrderCtrl

func GetBiananceCtrl() *BiananceOrderCtrl {
	if BianancePTR == nil {
		BianancePTR = new(BiananceOrderCtrl)
	}
	return BianancePTR
}

// 초기화 작업
func (ty *BiananceOrderCtrl) Init(acess, key string) {
	ty.B.Init(acess, key)
	ty.MarketList = make(map[string]futures.Symbol)
	ty.Leverage = make(map[string]futures.LeverageBracket)
}

// setClient
func (ty *BiananceOrderCtrl) SetUser(acess, key string) {
	ty.B.Init(acess, key)
}

// 마켓데이터 가져와서 저장하기 - 한번 부르고, 내부데이터 사용하면 됨
func (ty *BiananceOrderCtrl) GetMarketInfo() error {
	// 마켓데이터 가져오기
	res, err := ty.B.GetExchangeInfo()
	if err != nil {
		return err
	}

	// 마켓데이터 저장
	for _, v := range res.Symbols {
		ty.MarketList[v.Symbol] = v
	}
	return nil
}

// 레버리지 가져오기
func (ty *BiananceOrderCtrl) GetLeverageBracket() error {
	res, err := ty.B.GetLeverageList()
	if err != nil {
		return err
	}

	// 마켓데이터 저장
	for _, v := range res {
		ty.Leverage[v.Symbol] = *v
	}
	return nil
}

// 심볼의 최대 맥스 값 찾기
func (ty *BiananceOrderCtrl) GetMaxLeverage(symbol string) (futures.Bracket, error) {
	var tmp futures.Bracket
	leverageList := ty.Leverage[symbol]
	if len(leverageList.Brackets) == 0 {
		return tmp, fmt.Errorf("바스켓에 레버리지 없음")
	}
	return leverageList.Brackets[0], nil
}

// 바스켓 케파 확인하기 - 호출하면
func (ty *BiananceOrderCtrl) GetLeverageBracketDefault(symbol string, leverage int) (futures.Bracket, error) {
	/*

		레버리지 캐파는  50, 20, 10, 5 , 1  이렇게 나열되어 있고, 종목마다 다르다.
		레버리지 25 는 50 케파를 가지고 있다. 21~50까지는 50으로 적용한다.
		레버리지 20 는 11~20 까지 적용된다.

		레버리지 바스켓은 높은 순선대로 나열해서 준다 (바이낸스 기준)
		만약에 다른 거래소라면. 레버리지 값을 높은 순으로 정렬하는 구문을 추가해라.

		현재의 레버리지가 바스켓레버리지 보다 높으면, 이전 것을 준다.  25 라면 -> 50을 준다.
		현재의 레버리지가 바스켓레버리지 와 같다면 현재것을 준다. 20 -> 20 을준다.


		 Bracket          int     `json:"bracket"`          //     "bracket": 1,
		InitialLeverage  int     `json:"initialLeverage"`  //     "initialLeverage": 50,
		NotionalCap      float64 `json:"notionalCap"`      //     "notionalCap": 5000,
		NotionalFloor    float64 `json:"notionalFloor"`    //     "notionalFloor": 0,
		MaintMarginRatio float64 `json:"maintMarginRatio"` //     "maintMarginRatio": 0.015,
		Cum              float64 `json:"cum"`              //     "cum": 0


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

// 외부에서 사용하는 주문하기
func (ty *BiananceOrderCtrl) SendOrder(order BiananceOrder) (*futures.CreateOrderResponse, error) {

	// MarketList
	// Leverage

	// 마켓에 있는지 체크
	_, of := ty.MarketList[order.Symbol]
	if of == false {
		return nil, fmt.Errorf(NotFoundSymbo) // "없는 심볼"
	}

	// 유저 포지션 정보 ( GetPositionInfo(symbol string)

	// 최대 레버리지 가져와서 체크
	maxLeverage, err := ty.GetMaxLeverage(order.Symbol)
	if err != nil {
		return nil, fmt.Errorf("Binance API Error : %s", err)
	}

	// 현재 레버리지 및 포지션 정보
	userPosition, err := ty.B.GetPositionInfo(order.Symbol)
	if err != nil {
		return nil, fmt.Errorf("Binance API Error : %s", err)
	}

	// 등록되어 있는 현재 레버리지
	currentLeverage, err := strconv.Atoi(userPosition[0].Leverage)
	if err != nil {
		return nil, fmt.Errorf("Not Convert Lverage : %s", err)
	}

	// 현재 등록되어 있는건 의미가 있다? 없다?... 없다.. 계속 바꿀꺼니까.

	// 포지션이 있으면 레버리지 패스
	var havepositionflg = false
	for _, v := range userPosition {
		if v.EntryPrice != "0.0" {
			// return nil, fmt.Errorf("Have a Position")
			havepositionflg = true
		}
	}

	var nowLeverage int

	// 포지션이 없으면.  레버리지 따져야지.
	if havepositionflg {
		if order.Leverage <= maxLeverage.InitialLeverage {
			if order.Leverage != currentLeverage {
				ty.B.SetLeverage(order.Symbol, order.Leverage)
				nowLeverage = order.Leverage
			}
		} else {
			// 최대레버리지 보다 현재가 높으면 최대로 변경하고 진행
			ty.B.SetLeverage(order.Symbol, maxLeverage.InitialLeverage)
			nowLeverage = maxLeverage.InitialLeverage
		}
	}

	// 레버리지별 최대 수량 초과
	nowLeverageBarket, err := ty.GetLeverageBracketDefault(order.Symbol, nowLeverage)
	if err != nil {
		return nil, fmt.Errorf("Error GetLeverageBracketDefault :, %s", err.Error())
	}

	// 최소수량, 최소금액 만족하는 최소 수량 가져오기
	minqty, err := ty.GetMinQty(order)
	if err != nil {
		return nil, err
	}

	// 수량 변환
	quantity, err := utils.Float64(order.Quantity)
	if err != nil {
		return nil, err
	}

	// 최소바다 작으면 최소수량으로 맞춘다.
	if minqty > quantity {
		order.Quantity = utils.String(minqty)
	}

	// 최대금액 초과하는지 확인
	// 마진 체크하기 캐파가 되냐? 안되냐
	TotalAmount, err := ty.CheckMargin(order, nowLeverage)
	if err != nil {
		return nil, fmt.Errorf("Error CheckMargin :, %s", err.Error())
	}

	// 현재값이,  캐파보다 더 크다면.
	if TotalAmount > nowLeverageBarket.NotionalCap {
		order.Quantity = ty.ChangeQuantity(order, nowLeverageBarket.NotionalCap) // 케파에 맞게 줄이는 작업
	}

	if order.Type == "LIMIT" {
		// price 자릿수 계산하기
		// 가격 소수점 변경 - Makret은 할필요 없으니까.
		// PT/SL 은 저짝 따로 있고,
	}

	// 수량 소수점 - 자르는 것만 해주는게 낫겠네  그니까 대충 계산해서 넣으면 나오는거지.
	qty, err := ty.convertQty(order.Symbol, order.Quantity)
	if err != nil {
		return nil, fmt.Errorf("Error converqty :, %s", err.Error())
	}
	// 소수점 자릿수 변경해서 치환하기
	order.Quantity = qty

	res, err := ty.order(order)
	if err != nil {
		return nil, err
	}

	return res, err
}

// 내부용 주문하기
func (ty *BiananceOrderCtrl) order(order BiananceOrder) (*futures.CreateOrderResponse, error) {

	var orderInfo REST_FUTURE_NewOrder
	orderInfo.Symbol = order.Symbol
	orderInfo.Side = order.Side
	orderInfo.PositionSide = order.PositionSide
	orderInfo.OrderType = order.Type
	orderInfo.Quantity = order.Quantity

	send, err := ty.B.SendOrder(orderInfo)
	if err != nil {
		return nil, err
	}

	return send, nil
}

// TP 주문

// SL 주문

// 취소하기

// 수량 dot 구하기
func (ty *BiananceOrderCtrl) dotSizeForAmount(symbol string) (int, error) {
	coinInfo := ty.MarketList[symbol]

	// 0.001 => 3 으로 변경해야지?
	stepsize := coinInfo.MarketLotSizeFilter().StepSize

	b, cnt := utils.DecimalCount(stepsize)
	if b == false {
		return 0, fmt.Errorf("변환할 stepsize 에러 입니다. ")
	}

	return cnt, nil
}

// 수량을 최소 단위로 잘라주기 - 이거 안잘라주면 에러남
func (ty *BiananceOrderCtrl) convertQty(symbol string, qty string) (string, error) {

	cnt, err := ty.dotSizeForAmount(symbol)
	if err != nil {
		return "", err
	}

	amount, err := strconv.ParseFloat(qty, 64)
	if err != nil {
		return "", err
	}

	tmp := utils.RoundFloat(amount, cnt)

	send := utils.String(tmp)
	if err != nil {
		return "", err
	}

	return send, nil
}

// 최소수량, 최소금액을 만족하는 최소수량 가져오가 ㅣ
func (ty *BiananceOrderCtrl) GetMinQty(order BiananceOrder) (float64, error) {
	var send float64
	var price float64
	if order.Type == "LIMIT" {
		tmpprice, err := utils.Float64(order.Price)
		if err != nil {
			return -1, err
		}
		price = tmpprice
	} else {
		// []*futures.PriceChangeStats, error
		// res, err := ty.B.FutuClient.GetTickerInfo(order.Symbol)
		res, err := ty.B.GetTickerInfo(order.Symbol)
		if err != nil {
			return -1, err
		}

		for _, v := range res {
			tmp, err := utils.Float64(v.LastPrice)
			if err != nil {
				return -1, err
			}
			price = tmp
		}
	}

	curPrice := utils.String(price)
	symbolInfo := ty.MarketList[order.Symbol]
	minNotionalFIlter := symbolInfo.MinNotionalFilter()
	minMoney := minNotionalFIlter.Notional

	// 0.001 => 3 으로 변경해야지?
	minSize := symbolInfo.MarketLotSizeFilter().MinQuantity

	coinStep, err := ty.dotSizeForAmount(order.Symbol)
	if err != nil {
		return -1, err
	}

	send = utils.MinCoinSize(curPrice, minMoney, minSize, coinStep)

	return send, nil
}

func (ty *BiananceOrderCtrl) CheckMargin(order BiananceOrder, nowLeverage int) (float64, error) {

}

// 케파에 맞게 줄이는 작업
func (ty *BiananceOrderCtrl) ChangeQuantity(order BiananceOrder, NotionalCap float64) string {

}
