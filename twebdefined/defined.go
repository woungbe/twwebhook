package twebdefined

// 신규 주문 관련 설정
type REST_FUTURE_NewOrder struct {
	Symbol           string //	STRING	YES
	Side             string //	ENUM	YES
	PositionSide     string //	ENUM	NO	Default BOTH for One-way Mode ; LONG or SHORT for Hedge Mode. It must be sent in Hedge Mode.
	OrderType        string //	ENUM	YES
	TimeInForce      string //	ENUM	NO
	Quantity         string //	DECIMAL	YES	Cannot be sent with closePosition=true(Close-All)
	ReduceOnly       string //	STRING	NO	"true" or "false". default "false". Cannot be sent in Hedge Mode; cannot be sent with closePosition=true
	Price            string //	DECIMAL	NO
	NewClientOrderId string //	STRING	NO	A unique id among open orders. Automatically generated if not sent. Can only be string following the rule: ^[\.A-Z\:/a-z0-9_-]{1,36}$
	StopPrice        string //	DECIMAL	NO	Used with STOP/STOP_MARKET or TAKE_PROFIT/TAKE_PROFIT_MARKET orders.
	ClosePosition    string //	STRING	NO	true, false；Close-All，used with STOP_MARKET or TAKE_PROFIT_MARKET.
	ActivationPrice  string //	DECIMAL	NO	Used with TRAILING_STOP_MARKET orders, default as the latest price(supporting different workingType)
	CallbackRate     string //	DECIMAL	NO	Used with TRAILING_STOP_MARKET orders, min 0.1, max 5 where 1 for 1%
	WorkingType      string //	ENUM	NO	stopPrice triggered by: "MARK_PRICE", "CONTRACT_PRICE". Default "CONTRACT_PRICE"
	PriceProtect     string //	STRING	NO	"TRUE" or "FALSE", default "FALSE". Used with STOP/STOP_MARKET or TAKE_PROFIT/TAKE_PROFIT_MARKET orders.
	NewOrderRespType string //	ENUM	NO	"ACK", "RESULT", default "ACK"
}

type IBiananceOrder struct {
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
