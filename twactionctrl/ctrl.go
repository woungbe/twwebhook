package twactionctrl

import "fmt"

type ActionCtrl struct {
	// MarketInfo MarketInfo // 공용 마켓 정보
	Users UserInfos
}

var ActionCtrlPTR *ActionCtrl

func GetActionCtrl() *ActionCtrl {
	if ActionCtrlPTR == nil {
		ActionCtrlPTR = new(ActionCtrl)
		ActionCtrlPTR.Init()
	}
	return ActionCtrlPTR
}

// 초기화
func (ty *ActionCtrl) Init() {

	// 유저정보 가져오기
	// ty.UserInfo = new(ty.UserInfo)
	ty.GetUserData()

	// 마켓정보 가져오기
	GetMarektInfo().Init(key, val)

}

// 주문이 들어올때 처리하는 방법
func (ty *ActionCtrl) SendOrder(strategysrl int, symbol, side string) error {
	// 마켓정보
	err := ty.GetMarketInfo(symbol) // 마켓이 있는지 여부 확인 , 각종 사이즈 가져오기
	if err != nil {
		fmt.Println("마켓정보에 없는 데이터 호출했음")
		// return fmt.Println(err)
	}

	// 레버리지 바스켓 가져오기
	err := ty.GetLeverageBracket(symbol)
	if err != nil {
		return err
	}

	// 현재가 호출
	price := ty.GetPrice(symbol)
	for k, v := range ty.UserInfoList {
		if v.Checkstrategy(strategysrl) == true {
			v.NewOrder(v, symbol, side, price)
		}
	}

	return nil
}

// 주문하기
func (ty *ActionCtrl) NewOrder(userInfo, symbol, side string) error {
	// 그냥 주문 때려도 되는데.. 검색은 해야지

	TickSizeList := GetTickSizes()
	LeverageBracket := GetLeverageBracket()

	/*
		마켓데이터-틱사이즈 리스트
		 : 코인 자릿수 체크
		 : price 자릿수 체크
		 : 최소수량 예외처리
		 : 최대수량 예외처리
	*/

	// 최대레버리지와 현재 레버리지 비교 후 리턴 (레버리지 재정의)

	// 현재 레버리지 호출, 현재자산호출
	// 유저정보 가져오기, 현재가, ticksize, 레버리지, 보유자산,
	// := ty.GetUserData() // 주문 방식, 투입금액,
	amount := ty.GetUserData(price, TickSizeList, LeverageBracket) // 여기서 다 정리해서. 타입이랑, 주문수량 알려주기

	if userInfo.주문타입 == "1" {
		// 반대  side 청산 추가
	}

	// 진입금액이 레버리지 최대금액과 맞는지 확인하고  amount 리턴하기
	result := ty.Order()

	if userInfo.익절 == true {
		// 익절 추가
		tpPrice := userInfo.TP * result
		NewTP(symbol, side, tpPrice)
	}

	if userInfo.손절 == true {
		// 손절 추가
		tpPrice := userInfo.TP * result
		NewSL(symbol, side, tpPrice)
	}

	return nil
}

// TP // SL
func (ty *ActionCtrl) NewTP() {
	// 애는 상관없음 - 왜냐. 무조건 주문 들어가서 리턴 받은 다음 들어가니까.

	// 수량도 필요없고, 계산만 잘하면됨.
}

func (ty *ActionCtrl) NewSL() {
	// 애는 상관없음 - 왜냐. 무조건 주문 들어가서 리턴 받은 다음 들어가니까.

	// 수량도 필요없고, 계산만 잘하면됨.
}

// 유저정보 가져오기
func (ty *ActionCtrl) GetUserInfo(useridx int) {

	return
}
