package binanceorder

const (
	NotFoundSymbo = "없는 심볼"
)

type RESTErrorREQ struct {
	Code         int64  `json:"code"` //에러코드 OK = 성공
	ErrorMessage string `json:"msg"`  //에러메시지
}

// 주문, TP, SL 만 되는 주문 목록 - 주문간소화
type BiananceOrder struct {
	Symbol       string // BTCUSDT
	Side         string // BUY, SELL
	PositionSide string // LONG or SHORT
	Type         string // MARKET
	Quantity     string // 0.007
	Price        string // 수량계산할때의 값넣기
	StopPrice    string // 스탑 로스 확인
	WorkingType  string // 시장평균가: MARK_PRICE, 현재가:CONTRACT_PRICE
	Leverage     int    // 주문단계에서 레버리지 등록
}

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

// ExchangeInfo exchange info
type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []SymbolData  `json:"symbols"`
}

func (ty *ExchangeInfo) SetData() {

}

// RateLimit struct
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

type SymbolData struct {
	Symbol                string                   `json:"symbol"`
	Pair                  string                   `json:"pair"`
	ContractType          string                   `json:"contractType"`
	DeliveryDate          int64                    `json:"deliveryDate"`
	OnboardDate           int64                    `json:"onboardDate"`
	Status                string                   `json:"status"`
	MaintMarginPercent    string                   `json:"maintMarginPercent"`
	RequiredMarginPercent string                   `json:"requiredMarginPercent"`
	PricePrecision        int                      `json:"pricePrecision"`
	QuantityPrecision     int                      `json:"quantityPrecision"`
	BaseAssetPrecision    int                      `json:"baseAssetPrecision"`
	QuotePrecision        int                      `json:"quotePrecision"`
	UnderlyingType        string                   `json:"underlyingType"`
	UnderlyingSubType     []string                 `json:"underlyingSubType"`
	SettlePlan            int                      `json:"settlePlan"`
	TriggerProtect        string                   `json:"triggerProtect"`
	OrderType             []string                 `json:"OrderType"`
	TimeInForce           []string                 `json:"timeInForce"`
	Filters               []map[string]interface{} `json:"filters"`
	QuoteAsset            string                   `json:"quoteAsset"`
	MarginAsset           string                   `json:"marginAsset"`
	BaseAsset             string                   `json:"baseAsset"`
}

// LeverageBracket define the leverage bracket
type LeverageBracket struct {
	Symbol   string    `json:"symbol"`
	Brackets []Bracket `json:"brackets"`
}

// Bracket define the bracket
type Bracket struct {
	Bracket          int     `json:"bracket"`
	InitialLeverage  int     `json:"initialLeverage"`
	NotionalCap      float64 `json:"notionalCap"`
	NotionalFloor    float64 `json:"notionalFloor"`
	MaintMarginRatio float64 `json:"maintMarginRatio"`
	Cum              float64 `json:"cum"`
}
