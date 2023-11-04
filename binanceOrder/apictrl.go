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
	B BinanceInfo // 바이낸스 api 콜
}

// 초기화 작업
func (ty *BiananceOrderCtrl) Init(acess, key string) {
	ty.B.Init(acess, key)
}

// 외부에서 사용하는 주문하기
func (ty *BiananceOrderCtrl) SendOrder(order BiananceOrder) (*futures.CreateOrderResponse, error) {

	// MarketList
	// Leverage

	// 마켓에 있는지 체크
	if GetFutuExchangeInfo().IsTrading(order.Symbol) == false {
		return nil, fmt.Errorf(NotFoundSymbo) // "없는 심볼"
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
		if order.Leverage <= GetLeverageBraket().GetMaxLeverage(order.Symbol) {
			if order.Leverage != currentLeverage {
				ty.B.SetLeverage(order.Symbol, order.Leverage)
				nowLeverage = order.Leverage
			}
		} else {
			// 최대레버리지 보다 현재가 높으면 최대로 변경하고 진행
			ty.B.SetLeverage(order.Symbol, GetLeverageBraket().GetMaxLeverage(order.Symbol))
			nowLeverage = GetLeverageBraket().GetMaxLeverage(order.Symbol)
		}
	}

	// 레버리지별 최대 수량 초과
	NotionalCap, err := GetLeverageBraket().GetMaxAmount(order.Symbol, nowLeverage)
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
	if TotalAmount > NotionalCap {
		order.Quantity = ty.ChangeQuantity(order, NotionalCap) // 케파에 맞게 줄이는 작업
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
func (ty *BiananceOrderCtrl) sendTP(order BiananceOrder, TPPerSent string) (*futures.CreateOrderResponse, error) {
	return ty.SendOrderTPSL(order, "TAKE_PROFIT_MARKET", TPPerSent)
}

// SL 주문
func (ty *BiananceOrderCtrl) sendSL(order BiananceOrder, SLPerSent string) (*futures.CreateOrderResponse, error) {
	return ty.SendOrderTPSL(order, "STOP_MARKET", SLPerSent)
}

func (ty *BiananceOrderCtrl) SendOrderTPSL(order BiananceOrder, types string, TPSLPrice string) (*futures.CreateOrderResponse, error) {
	var err error
	var orderSL REST_FUTURE_NewOrder
	tmpSide := ""

	// 익절이든 손절이든.
	if order.Side == "BUY" {
		tmpSide = "SELL"
	} else {
		tmpSide = "BUY"
	}

	orderSL.Symbol = order.Symbol
	orderSL.Side = tmpSide // order.Side 반대로가.
	orderSL.PositionSide = order.PositionSide
	orderSL.OrderType = types
	orderSL.ClosePosition = "true"
	orderSL.WorkingType = order.WorkingType
	// fmt.Println("orderSL : ", orderSL)
	var send *futures.CreateOrderResponse

	if types == "STOP_MARKET" {
		orderSL.StopPrice, err = ty.CurrentStopPrice(order.Symbol, orderSL.PositionSide, strconv.Itoa(order.Leverage), TPSLPrice)
	}

	if types == "TAKE_PROFIT_MARKET" {
		orderSL.StopPrice, err = ty.CurrentProfitPrice(order.Symbol, orderSL.PositionSide, strconv.Itoa(order.Leverage), TPSLPrice)
	}

	if err != nil {
		fmt.Println("익절 손절 주문에 실패했습니다. ")
	}

	tmpsend, er := ty.B.SendOrder(orderSL)
	if er != nil {
		fmt.Println("orderSL error : ", er)
	}
	send = tmpsend

	return send, nil
}

func (ty *BiananceOrderCtrl) CurrentStopPrice(symbol, side, Leverage, present string) (string, error) {
	var send string

	lastPrice, err := ty.GetTickerInfo(symbol)
	if err != nil {
		return "", err
	}
	m_price, _ := utils.Float64(lastPrice)
	m_leverage, _ := utils.Float64(Leverage)
	m_persent, _ := utils.Float64(present)

	priceLotSize := GetFutuExchangeInfo().GetExchangeLotTickInfo(symbol).GettickDotSize()

	if side == "LONG" {
		tmp := utils.ToFixed(((1 - (m_persent / m_leverage)) * m_price), priceLotSize)
		send = utils.String(tmp)
	} else if side == "SHORT" {
		tmp := utils.ToFixed(((1 + (m_persent / m_leverage)) * m_price), priceLotSize)
		send = utils.String(tmp)
	}

	return send, nil
}

func (ty *BiananceOrderCtrl) CurrentProfitPrice(symbol, side, Leverage, present string) (string, error) {
	var send string

	// 현재가 가져오기
	lastPrice, err := ty.GetTickerInfo(symbol)
	if err != nil {
		return "", err
	}

	// 현재가 ,레버리지, 퍼센트, priceLotSize
	m_price, _ := utils.Float64(lastPrice)
	m_leverage, _ := utils.Float64(Leverage)
	m_persent, _ := utils.Float64(present)
	priceLotSize := GetFutuExchangeInfo().GetExchangeLotTickInfo(symbol).GettickDotSize()

	if side == "LONG" {
		tmp := utils.ToFixed(((1 + (m_persent / m_leverage)) * m_price), priceLotSize)
		send = utils.String(tmp)
	} else if side == "SHORT" {
		tmp := utils.ToFixed(((1 - (m_persent / m_leverage)) * m_price), priceLotSize)
		send = utils.String(tmp)
	}

	return send, nil
}

// 현재가 가져오기
func (ty *BiananceOrderCtrl) GetTickerInfo(symbol string) (string, error) {
	ticker, err := ty.B.GetTickerInfo(symbol)
	if err != nil {
		return "", err
	}

	for _, v := range ticker {
		if v.Symbol == symbol {
			return v.LastPrice, nil
		}
	}

	return "", nil
}

// 취소하기

// 수량을 최소 단위로 잘라주기 - 이거 안잘라주면 에러남
func (ty *BiananceOrderCtrl) convertQty(symbol string, qty string) (string, error) {
	cnt := GetFutuExchangeInfo().GetExchangeLotTickInfo(symbol).GettickDotSize()
	amount, err := strconv.ParseFloat(qty, 64)
	if err != nil {
		return "", err
	}

	// 반올림 아니고, 잘라야됨 !!
	tmp := utils.CeilFloat(amount, cnt)
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
	minMoney := GetFutuExchangeInfo().GetExchangeLotTickInfo(order.Symbol).GetMinNotional()
	// 0.001 => 3 으로 변경해야지?
	minSize := GetFutuExchangeInfo().GetExchangeLotTickInfo(order.Symbol).GetMinQuantity()
	coinStep := GetFutuExchangeInfo().GetExchangeLotTickInfo(order.Symbol).GettickDotSize()
	send = utils.MinCoinSize(curPrice, minMoney, minSize, coinStep)
	return send, nil
}

// 마진 타입 맞는지
func (ty *BiananceOrderCtrl) CheckMargin(order BiananceOrder, nowLeverage int) (float64, error) {

	return 0, nil
}

// 케파에 맞게 줄이는 작업
func (ty *BiananceOrderCtrl) ChangeQuantity(order BiananceOrder, NotionalCap float64) string {
	return "0"
}
