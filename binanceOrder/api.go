package binanceorder

import (
	"context"
	"fmt"
	"twwebhook/twebdefined"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

/*
여려명을 해야 되기 때문에 userInfo 의 하위로 빠지게 됨 .
*/

type BinanceInfo struct {
	Client     *binance.Client
	FutuClient *futures.Client
}

func (ty *BinanceInfo) Init(key1, key2 string) {
	ty.Client = binance.NewClient(key1, key2)
	ty.FutuClient = binance.NewFuturesClient(key1, key2) // == futures.NewClient(apiKey, secretKey)
}

// 마켓데이터 가져오기
func (ty *BinanceInfo) GetExchangeInfo() (*futures.ExchangeInfo, error) {
	exInfo, err := ty.FutuClient.NewExchangeInfoService().Do(context.Background())
	// times, err := ty.Client.NewServerTimeService().Do(context.Background())
	if err != nil {
		fmt.Println("그냥 해보라고요? : ", err)
		return nil, err
	}

	return exInfo, nil
}

// 바이낸스 시간 데이터 가져오기
func (ty *BinanceInfo) GetServiceTime() (int64, error) {
	times, err := ty.Client.NewServerTimeService().Do(context.Background())
	if err != nil {
		fmt.Println("그냥 해보라고요? : ", err)
		return 0, err
	}
	return times, nil
}

// 바이낸스 펀딩피 리스트 가져오기
func (ty *BinanceInfo) GetFunddingFeeList() ([]*futures.PremiumIndex, error) {
	// res []*futures.PremiumIndex, err error
	res, err := ty.FutuClient.NewPremiumIndexService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 레버리지 바스켓 가져오기
func (ty *BinanceInfo) GetLeverageList() ([]*futures.LeverageBracket, error) {
	// GetLeverageBracketService
	res, err := ty.FutuClient.NewGetLeverageBracketService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 레버리지 세팅하기
func (ty *BinanceInfo) SetLeverage(symbol string, leverage int) (*futures.SymbolLeverage, error) {
	// ChangeLeverageService symbol  leverage
	res, err := ty.FutuClient.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 주문 넣기
func (ty *BinanceInfo) SendOrder(inf twebdefined.REST_FUTURE_NewOrder) (*futures.CreateOrderResponse, error) {
	// CreateOrderService
	obj := ty.FutuClient.NewCreateOrderService().
		Symbol(inf.Symbol).
		Side(futures.SideType(inf.Side)).
		Type(futures.OrderType(inf.OrderType)).
		Quantity(inf.Quantity)

	if inf.NewOrderRespType == "" {
		obj.NewOrderResponseType("ACK")
	} else {
		obj.NewOrderResponseType(futures.NewOrderRespType(inf.NewOrderRespType))
	}

	if inf.PositionSide != "" {
		obj.PositionSide(futures.PositionSideType(inf.PositionSide))
	}
	if inf.TimeInForce != "" {
		obj.TimeInForce(futures.TimeInForceType(inf.TimeInForce))
	}
	if inf.ReduceOnly != "" {
		var v bool
		v = false
		if inf.ReduceOnly == "true" {
			v = true
		}
		obj.ReduceOnly(v)
	}
	if inf.Price != "" {
		obj.Price(inf.Price)
	}
	if inf.NewClientOrderId != "" {
		obj.NewClientOrderID(inf.NewClientOrderId)
	}
	if inf.StopPrice != "" {
		obj.StopPrice(inf.StopPrice)
	}

	if inf.ClosePosition != "" {
		var v bool
		v = false
		if inf.ClosePosition == "true" {
			v = true
		}
		obj.ClosePosition(v)
	}
	if inf.ActivationPrice != "" {
		obj.ActivationPrice(inf.ActivationPrice)
	}
	if inf.CallbackRate != "" {
		obj.CallbackRate(inf.CallbackRate)
	}
	if inf.WorkingType != "" {
		obj.WorkingType(futures.WorkingType(inf.WorkingType))
	}
	if inf.PriceProtect != "" {
		var v bool
		v = false
		if inf.ClosePosition == "TRUE" {
			v = true
		}
		obj.PriceProtect(v)
	}

	res, er := obj.Do(context.Background())
	fmt.Println("주문 정보 : ", inf)
	fmt.Println("주문 리턴 : ", res)
	if er != nil {
		return nil, er
	}

	return res, nil
}

// 자산 조회
func (ty *BinanceInfo) GetAccountInfo() ([]*futures.Balance, error) {
	// GetBalanceService
	res, err := ty.FutuClient.NewGetBalanceService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 현재가 가져오기  (ticker )
func (ty *BinanceInfo) GetTickerInfo() ([]*futures.PriceChangeStats, error) {
	// /v1/ticker/24hr - ListPriceChangeStatsService
	res, err := ty.FutuClient.NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 현재 포지션 가져오기 정보 -- 펀딩피에서는 사용안함
func (ty *BinanceInfo) GetPositionInfo() ([]*futures.PositionRisk, error) {
	// /v2/positionRisk - GetPositionRiskService - GetPositionRiskService
	res, err := ty.FutuClient.NewGetPositionRiskService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 미체결 모두 날리기 !!
func (ty *BinanceInfo) SetAllOpenOrderCancel(symbol string) error {
	// /v2/positionRisk - GetPositionRiskService - GetPositionRiskService
	err := ty.FutuClient.NewCancelAllOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// fapi/v1/allOpenOrders

// marginType 가져오기 -- ISOLATED, CROSSED
func (ty *BinanceInfo) SetMarginType(symbol string, marginType futures.MarginType) error {
	// /v2/positionRisk - GetPositionRiskService - GetPositionRiskService
	err := ty.FutuClient.NewChangeMarginTypeService().Symbol(symbol).MarginType(marginType).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
