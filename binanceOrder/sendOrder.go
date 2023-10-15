package binanceorder

type SendOrder struct {
	C     *BinanceInfo
	Order BiananceOrder
}

type BiananceOrder struct {
	Symbol          string // BTCUSDT
	Side            string // BUY, SELL
	PositionSide    string // LONG or SHORT
	Type            string // MARKET
	Quantity        string // 0.007
	Leverage        string // 레버리지
	LastFunddingFee string // 펀딩피 체크할 때 사용함
	Price           string // 수량계산할때의 값넣기
	StopPrice       string // 스탑 로스 확인
	WorkingType     string // 시장평균가: MARK_PRICE, 현재가:CONTRACT_PRICE
}
