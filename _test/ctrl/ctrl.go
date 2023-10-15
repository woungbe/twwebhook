package ctrl

import binanceorder "twwebhook/binanceOrder"

/*
신호가 오면 주문을 신청하는 로직
*/

type BianaceSendOrderCtrl struct {
	// 마켓정보
	ExchangeInfo BianaceExchangeInfo
	UserInfo     UserInfo
	SignalList   TwSignal
	// 유저정보 데이터, key, val 도 여기에 있음
	// 신호 리스트 - 이게 맞나 ?...
}

type OrderDataInfo struct {
	UserIDX string // 오더 정보를 완성하려면 유저 정보가 필요함
}

var orderCtrlPTR *BianaceSendOrderCtrl

// 초기화 시작
func Init() {
	orderCtrlPTR = new(BianaceSendOrderCtrl)

	// 마켓정보
	// 유저정보
	// 신호리스트 - 정보 가져오기

}

// 데이터 가져오기
func GetCtrl() *BianaceSendOrderCtrl {
	if orderCtrlPTR == nil {
		Init()
	}

	return orderCtrlPTR
}

// 주문전 주문할 데이터 만들기 -
func (ty *BianaceSendOrderCtrl) beforeOrderInfo(data OrderDataInfo) (binanceorder.BiananceOrder, error) {

	// 데이터로 받을 목록
	// symbol
	// leverage
	// price
	// long, short

	// DB로 받을 데이터
	// marginType
	// 익절 비율
	// 손절 비율
	// 시장가, 지정가
	// postOnly

	// 목표는 틀리냐 마냐가 아니라.
	// BiananceOrder를 만드는 것
	// 웬만한 예외처리는 여기서 다 처리합시다.

	ty.checkMargin(send)
	ty.checkLeverage(send)
	ty.checkUserAssets(send)
	ty.checkPriceDecimal(send)

	// 가격 소수점 체크, 수량 소수점 체크
	price, err := checkPriceDecimal(send)
	if err != nil {
		return err
	}

	var send binanceorder.BiananceOrder
	return send, nil

}

// 주문하기
func (ty *BianaceSendOrderCtrl) sendOrder(data OrderDataInfo) error {

	// 마진타입이 맞는지 체크
	// MarginType == send
	// 레버리지 체크
	// 유저금액이 있는지 체크

	var send binanceorder.BiananceOrder
	send, err := beforeOrderInfo(data)
	if err != nil {
		return err
	}

	// 데이터 검증
	if CheckSendOrderReturn(data, send) {
		newTP(data, send) // 리턴 없이 내부처리 시킴
		newSL(data, send) // 리턴 없이 내부처리 시킴
	}

	return nil
}

// 마진 타입 체크 필요
func (ty *BianaceSendOrderCtrl) checkMargin(send binanceorder.BiananceOrder) error {
	// 마진타입가져와서 맞는지 체크해야됨
	// 아니면 현재 있는 마진 타입으로 진입시켜야됨.

	// 여기서 변경요청하면 됨
	return nil
}

// 레버리지 확인 필요 - 최대 레버리지 초과면 최대로 변경
func (ty *BianaceSendOrderCtrl) checkLeverage(send binanceorder.BiananceOrder) error {

	// 여기서 변경요청 하면됨 .
	return nil
}

// 가격과 수량을 체크해주는 것 필요
func (ty *BianaceSendOrderCtrl) checkUserAssets(send binanceorder.BiananceOrder) (string, error) {
	// 유저정보를 토대로 얼마나 사겠다가 필요함
	return "", nil
}

// 가격 소수점 자르기
func (ty *BianaceSendOrderCtrl) checkPriceDecimal(send binanceorder.BiananceOrder) (string, error) {
	// 마켓정보에서 소수점 찾아서 준비시킴
	// 해당 소수점으로 정리해서 보내줄 필요가 있음
	return "", nil
}

// sendOrderReturn
func (ty *BianaceSendOrderCtrl) CheckSendOrderReturn(send binanceorder.BiananceOrder) bool {

	return false
}

// TP
func (ty *BianaceSendOrderCtrl) newTP(data OrderDataInfo, send binanceorder.BiananceOrder) {
	// 유저 익절데이터 있고,

	// 익절요청하고 실패하면 2배 재요청

}

// SL
func (ty *BianaceSendOrderCtrl) newSL(data OrderDataInfo, send binanceorder.BiananceOrder) {

	// 손절 요청하고 실패하면
	// 강제 청산 또는 메시지 전송

}
