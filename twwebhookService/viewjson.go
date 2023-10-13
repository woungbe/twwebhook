package twwebhookService

import (
	"fmt"
	"strconv"
	"strings"
)

type REST_WEB_TradingViewJson struct {
	StrategyName string `json:"strategyName"`
	Coin         string `json:"coin"`
	PositionSide string `json:"positionSide"`
	PowerText    string `json:"powerText"`
	Leverage     string `json:"leverage"`
	MarginType   string `json:"marginType"`
	Price        string `json:"price"`
	Profit       string `json:"profit"`
	Losscut      string `json:"losscut"`
	BottomText   string `json:"bottomText"`
}

func (ty REST_WEB_TradingViewJson) MakeText() string {
	Profit := ty.GetProfitText()
	Losscut := ty.GetLoseCutText()

	send := fmt.Sprintf(`%s
종목 : %s,
포지션 : %s%s
레버리지 : %sx%s,
진입가 : %s,
익절 : %s
손절 : %s

%s
`,
		ty.StrategyName,
		ty.Coin,
		ty.GetChangeTextSide(),
		ty.PowerText,
		ty.Leverage,
		ty.MarginType,
		ty.Price,
		Profit,
		Losscut,
		ty.BottomText,
	)

	/*
		슈퍼트레이더 트레이딩 시그널

		종목 : {{ticker}},
		포지션 : 숏(강력),
		레버리지 : 10x(격리),
		가격 : {{close}},
		익절 : 0.97
		손절 : 1.03

		★항상, 매번 말씀드리지만, 꼭 손절을 지키면서 매매 진행하시길 바랍니다.
		( 본인 투자에 대한 책임은 본인에게 있습니다)
		💡누적 수익으로 항상 접근하셔서 뇌동매매를 방지하시길 바랍니다!
	*/
	return send
}

func (ty REST_WEB_TradingViewJson) GetChangeTextSide() string {
	retText := ""
	if ty.PositionSide == "LONG" {
		retText = "롱"
	} else if ty.PositionSide == "SHORT" {
		retText = "숏"
	}
	return retText
}

func (ty REST_WEB_TradingViewJson) GetProfitText() string {
	retText := ""

	// 가격
	price, err := strconv.ParseFloat(ty.Price, 64)
	if err != nil {
		return err.Error()
	}

	// 익절
	profit, err := strconv.ParseFloat(ty.Profit, 64)
	if err != nil {
		return err.Error()
	}

	b, decimalPlaces := DecimalCount(ty.Price)
	if b == false {
		return ""
	}

	retText = strconv.FormatFloat((price * profit), 'f', decimalPlaces, 64)
	return retText
}

func (ty REST_WEB_TradingViewJson) GetLoseCutText() string {
	retText := ""
	// 가격
	price, err := strconv.ParseFloat(ty.Price, 64)
	if err != nil {
		return err.Error()
	}

	// 익절
	losscut, err := strconv.ParseFloat(ty.Losscut, 64)
	if err != nil {
		return err.Error()
	}

	b, decimalPlaces := DecimalCount(ty.Price)
	if b == false {
		return ""
	}

	retText = strconv.FormatFloat((price * losscut), 'f', decimalPlaces, 64)
	return retText

}

// 소수점 몇째자리인가 카운트 하는것
func DecimalCount(price string) (bool, int) {
	value, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println("부동소수점 변환 오류:", err)
		return false, -1
	}

	var decimalPlaces int
	decimalPlaces = -1

	// 부동소수점 값을 문자열로 변환하여 소수점 이하 자릿수 계산
	strValue := strconv.FormatFloat(value, 'f', -1, 64)
	parts := strings.Split(strValue, ".")
	if len(parts) == 2 {
		decimalPlaces = len(parts[1])
		fmt.Println("소수점 이하 자릿수:", decimalPlaces)
	} else {
		if len(parts[0]) == 1 {
			i, _ := strconv.Atoi(parts[0])
			if i < 10 {
				return true, 0
			}
		}
		fmt.Println("부동소수점 형식이 아닙니다.")
		return false, decimalPlaces
	}

	return true, decimalPlaces
}

/*
	body, _ := io.ReadAll(r.Body)
	//bodyString := string(body)

	err := json.Unmarshal(body, &inf) // .NewDecoder(r.Body).Decode(&inf)
	if err != nil {
		reqError.ErrorCode = "SystemError_0101"
		reqError.ErrorMessage = "body data error"

		jsondata, _ := json.Marshal(reqError)
		http.Error(w, string(jsondata), http.StatusBadRequest)
		return
	}
*/
