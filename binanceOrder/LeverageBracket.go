package binanceorder

import (
	"context"
	"fmt"
	"twwebhook/utils"

	"github.com/adshao/go-binance/v2/futures"
)

type LeverageBraket struct {
	MapLeverageBracket map[string]*futures.LeverageBracket
	LastUpdateTime     int64 //마지막 업데이트 시간
}

/*
type Bracket struct {
	Bracket          int     `json:"bracket"`
	InitialLeverage  int     `json:"initialLeverage"`
	NotionalCap      float64 `json:"notionalCap"`
	NotionalFloor    float64 `json:"notionalFloor"`
	MaintMarginRatio float64 `json:"maintMarginRatio"`
	Cum              float64 `json:"cum"`
}
*/

var g_LeverageBraketPTR *LeverageBraket

func GetLeverageBraket() *LeverageBraket {
	if g_LeverageBraketPTR == nil {
		g_LeverageBraketPTR = new(LeverageBraket)
	}
	return g_LeverageBraketPTR
}

// 초기화
func (ty *LeverageBraket) Init() {
	ty.MapLeverageBracket = make(map[string]*futures.LeverageBracket)
	ty.LastUpdateTime = 0
}

// 바이낸스에서 데이터 받기
func (ty *LeverageBraket) GetLeverageBracket() error {
	bin := futures.NewClient("", "")
	res, err := bin.NewGetLeverageBracketService().Do(context.Background())
	if err != nil {
		return err
	}

	for _, v := range res {
		ty.MapLeverageBracket[v.Symbol] = v
	}

	return nil
}

// 레버리지정보 준비되어있는가
func (ty *LeverageBraket) IsRady() (isOk bool, reload bool) {
	defer func() {
		if err := recover(); err != nil {
			// pawlog.Error("Crit Panic", "Error", err)
		}
	}()

	curTime := utils.GetCurrentTimestamp()
	if len(ty.MapLeverageBracket) <= 0 {
		return false, true
	}
	//== 여기는 재호출 할것인가.
	if (ty.LastUpdateTime + (60 * 60)) < curTime {
		ty.LastUpdateTime = curTime - (60 * 59)
		return true, true
	}
	return true, false
}

//===========================================

// 바스켓 케파 찾기
func (ty *LeverageBraket) GetBarcket(symbol string, leverage int) (futures.Bracket, error) {
	defer func() {
		if err := recover(); err != nil {
			// pawlog.Error("Crit Panic", "Error", err)
		}
	}()

	leverageList := ty.MapLeverageBracket[symbol]
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

//================================

// 현재 (심볼,레버리지)  최대 가능 금액 리턴해주기
func (ty *LeverageBraket) GetMaxAmount(symbol string, leverage int) (float64, error) {
	defer func() {
		if err := recover(); err != nil {
			// pawlog.Error("Crit Panic", "Error", err)
		}
	}()

	send, err := ty.GetBarcket(symbol, leverage)
	return send.NotionalCap, err
}

// 해당 심볼의 최대 레버리지를 가져온다
func (ty *LeverageBraket) GetMaxLeverage(symbol string) int {
	defer func() {
		if err := recover(); err != nil {
			// pawlog.Error("Crit Panic", "Error", err)
		}
	}()

	obj, v := ty.MapLeverageBracket[symbol]
	if !v {
		return 1
	}

	var maxlv int
	maxlv = 1
	for i := 0; i < len(obj.Brackets); i++ {
		if obj.Brackets[i].InitialLeverage > maxlv {
			maxlv = obj.Brackets[i].InitialLeverage
		}
	}
	return maxlv
}
