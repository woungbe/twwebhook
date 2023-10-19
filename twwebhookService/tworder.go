package twwebhookService

// 데이터 가져오기
type REST_NewOrder struct {
	Accesskey     string `json:"accesskey"`
	Secritkey     string `json:"secritkey"`
	Price         string `json:"price"`
	MarginType    string `json:""`             // CORSS, ISOLATED
	Types         string `json:"types"`        // MARKET, LIMIT
	Side          string `json:"side"`         // LONG, SHORT
	PositionSide  string `json:"positionSide"` // BUY,SELL
	Profit        string `json:"profit"`       // 1.03
	Losscut       string `json:"losscut"`      // 0.97
	AmountPersent string `json:"amount"`       // 주문 수량
}

// 시나리오 정리하기
func (ty REST_NewOrder) SendOrder() {
	// 주문 타입으로 만들어야되구요

	// 익절 손절 계산해야되구요

	// 익절 손절에 맞게 주문 넣어야되구요 .

	// MarginType 변경 체크
	// 최대 레버리지 체크 및 변경

	// 수량 소수점 정리

	// 시나리오 좀 크네 .. 따로 만들자. 나도 쓰게.
	// 그럼 DB도 좀 넣구.. 이것저것 추가하자.
	// 컨트롤 하게 사이즈도 키우면 괜찮을듯.
}
