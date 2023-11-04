package binanceorder

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"twwebhook/twwebhookConfig"

	"github.com/adshao/go-binance/v2/futures"
)

/*
	마켓 정보 관리
*/

type ExchangeInfoData struct {
	mbInit          bool
	mExchangeData   *futures.ExchangeInfo
	mMapLotTickInfo map[string]*ExchangeLotTickInfo
	mMapSymbols     map[string]futures.Symbol
}

type ExchangeLotTickInfo struct {
	tickSize    string  //가격 틱 사이즈  BTc =>0.1
	fTickSize   float64 //가격 틱사이즈 실수로 => 0.1
	tickDotSize int     // 틱사이즈의 소숫점 이하 자릿수 => 1
	lotSize     string  //롯 사이즈 원문자 ( 주문시 수량 단위) => 0.001 수량단위
	fLotSize    float64 //롯사이즈 실수로 => 0.001
	lotDotSize  int     //롯사이즈의 소숫점 이하 자릿수 => 3

	MinQuantity  string  //최소 거래 수량 => 0.001
	fMinLotSize  float64 //최소 거래 수량 => 0.001
	MinNotional  string  //최소 거래 금액 => 5
	fMinNotional float64 //최소 거래 금액 => 5

	MarketMaxQuantity string  //시장가 최대 주문 가능 수량
	fMarketMaxLotSize float64 //시장가 최대 주문 가능 수량
}

var g_ExchangeLotTickInfoPTR *ExchangeInfoData

func GetFutuExchangeInfo() *ExchangeInfoData {
	if g_ExchangeLotTickInfoPTR == nil {
		g_ExchangeLotTickInfoPTR = new(ExchangeInfoData)
		//g_ExchangeLotTickInfoPTR.Init()
	}
	return g_ExchangeLotTickInfoPTR
}

// 초기화
func (ty *ExchangeInfoData) Init() {
	//ty.mExchangeData = new(ExchangeInfo)
	if ty.mbInit {
		fmt.Println("ExchangeInfoData 이미초기화됌")
		return
	}
	ty.mbInit = true
	ty.mMapLotTickInfo = make(map[string]*ExchangeLotTickInfo)
	ty.mMapSymbols = make(map[string]futures.Symbol)
	ty.reLoadExchangeInfo()
	go ty.goSchExchangeInfo()
}

// 틱롯 정보 객체 리턴
func (ty *ExchangeInfoData) GetExchangeLotTickInfo(symbol string) *ExchangeLotTickInfo {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()
	obj := ty.mMapLotTickInfo[symbol]
	return obj
}

// 거래중인가 체크 - 존재하는 것도 같이 체크
func (ty *ExchangeInfoData) IsTrading(symbol string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	obj, v := ty.mMapSymbols[symbol]
	if !v {
		return false
	}
	if obj.Status != "TRADING" {
		return false
	}
	return true

}

// 하루 한번 정보 업뎃
func (ty *ExchangeInfoData) goSchExchangeInfo() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	for {
		if twwebhookConfig.IsDebugmode() {
			time.Sleep(time.Minute * 3)
		} else {
			time.Sleep(time.Hour)
		}

		ty.reLoadExchangeInfo()
	}
}

// 바이넨스에서 정보 가져온다.
func (ty *ExchangeInfoData) reLoadExchangeInfo() bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	bin := futures.NewClient("", "")
	receiveData, er := bin.NewExchangeInfoService().Do(context.Background())
	if er != nil {
		fmt.Println("Error recover ", er)
		return false
	}

	// var newData ExchangeInfo
	// newData.SetData(receiveData)
	ty.mExchangeData = receiveData

	for _, sy := range ty.mExchangeData.Symbols {
		ty.mMapSymbols[sy.Symbol] = sy
		org, v := ty.mMapLotTickInfo[sy.Symbol]
		if !v {
			pf := sy.PriceFilter()
			lf := sy.LotSizeFilter()
			mi := sy.MinNotionalFilter()
			marketLot := sy.MarketLotSizeFilter()

			newExp := new(ExchangeLotTickInfo)
			newExp.tickSize = pf.TickSize
			newExp.lotSize = lf.StepSize

			newExp.MarketMaxQuantity = marketLot.MaxQuantity

			newExp.MinNotional = "0"
			if mi != nil {
				newExp.MinNotional = mi.Notional
			}

			newExp.MinQuantity = lf.MinQuantity

			newExp.createBasicInfo()
			ty.mMapLotTickInfo[sy.Symbol] = newExp
		} else {
			pf := sy.PriceFilter()
			lf := sy.LotSizeFilter()
			mi := sy.MinNotionalFilter()
			marketLot := sy.MarketLotSizeFilter()

			org.tickSize = pf.TickSize
			org.lotSize = lf.StepSize

			org.MarketMaxQuantity = marketLot.MaxQuantity

			org.MinNotional = "0"
			if mi != nil {
				org.MinNotional = mi.Notional
			}

			org.MinQuantity = lf.MinQuantity

			org.createBasicInfo()
		}
	}
	return true
}

// GetMinNational 최소 주문 금액
func (ty *ExchangeInfoData) GetMinNational(symbol string) float64 {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	obj := ty.GetExchangeLotTickInfo(symbol)
	if obj == nil {
		return 0
	}
	return obj.fMinNotional
}

// GetMinLotSize 최소 거래 수량
func (ty *ExchangeInfoData) GetMinLotSize(symbol string) float64 {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()
	obj := ty.GetExchangeLotTickInfo(symbol)
	if obj == nil {
		return 0
	}
	return obj.fMinLotSize
}

// GetMaxMarketQuantity 시장가 최대 주문 가능 수량
func (ty *ExchangeInfoData) GetMaxMarketQuantity(symbol string) float64 {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()
	obj := ty.GetExchangeLotTickInfo(symbol)
	if obj == nil {
		return 0
	}
	return obj.fMarketMaxLotSize
}

//=================================

// 틱사이즈 문자
func (ty *ExchangeLotTickInfo) GetTickSize_Str() string {
	return ty.tickSize
}

func (ty *ExchangeLotTickInfo) GetTickSize_F() float64 {
	return ty.fTickSize
}

func (ty *ExchangeLotTickInfo) GetLotSize_Str() string {
	return ty.lotSize
}

func (ty *ExchangeLotTickInfo) GetLotSize_F() float64 {
	return ty.fLotSize
}

func (ty *ExchangeLotTickInfo) GetMaxMarketQuantity_F() float64 {
	return ty.fMarketMaxLotSize
}

func (ty *ExchangeLotTickInfo) GetMinLotSize_F() float64 {
	return ty.fMinLotSize
}

func (ty *ExchangeLotTickInfo) GetMinNotionalSize_F() float64 {
	return ty.fMinNotional
}

func (ty *ExchangeLotTickInfo) GettickDotSize() int {
	return ty.tickDotSize
}

func (ty *ExchangeLotTickInfo) GetLotDotSize() int {
	return ty.lotDotSize
}

func (ty *ExchangeLotTickInfo) GetMinQuantity() string {
	return ty.MinQuantity
}

func (ty *ExchangeLotTickInfo) GetfMinLotSize() float64 {
	return ty.fMinLotSize
}

func (ty *ExchangeLotTickInfo) GetMinNotional() string {
	return ty.MinNotional
}

func (ty *ExchangeLotTickInfo) GetfMinNotional() float64 {
	return ty.fMinNotional
}

// 기본 정보로 기타 자료 생성
func (ty *ExchangeLotTickInfo) createBasicInfo() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()
	ty.fTickSize, _ = strconv.ParseFloat(ty.tickSize, 64)
	ty.fLotSize, _ = strconv.ParseFloat(ty.lotSize, 64)
	ty.tickDotSize = ty.getFloatDotSize(ty.tickSize)
	ty.lotDotSize = ty.getFloatDotSize(ty.lotSize)
	ty.fMinLotSize, _ = strconv.ParseFloat(ty.MinQuantity, 64)
	ty.fMinNotional, _ = strconv.ParseFloat(ty.MinNotional, 64)
	ty.fMarketMaxLotSize, _ = strconv.ParseFloat(ty.MarketMaxQuantity, 64) //시장가 최대 주문 가능 수량
}

// 실수의 경우 뒤 0을 제외한 소숫점 자리수를 리턴 (문자로 줘야함)
func (ty *ExchangeLotTickInfo) getFloatDotSize(fData string) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	l := len(fData)
	dot := false
	ss := 0
	for i := 0; i < l; i++ {
		if fData[i] == '.' {
			dot = true
		} else if dot && fData[i] != '0' {
			ss++
			return ss
		} else if dot && fData[i] == '0' {
			ss++
		}
	}
	return 0
}

// 틱사이즈= 지정 소숫점하위를 잘라서 리턴
func (ty *ExchangeLotTickInfo) ConvetTickDotSize(fData string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	if len(fData) <= 0 {
		return fData
	}

	return ty.convetDotSize(fData, ty.tickDotSize)
}

// 롯사이즈= 지정 소숫점하위를 잘라서 리턴
func (ty *ExchangeLotTickInfo) ConvetLotDotSize(fData string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	if len(fData) <= 0 {
		return fData
	}
	return ty.convetDotSize(fData, ty.lotDotSize)
}

// 소숫점 하위 잘라서 리턴
func (ty *ExchangeLotTickInfo) convetDotSize(fData string, dotSize int) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	if len(fData) <= 0 {
		return fData
	}

	arr := strings.Split(fData, ".")

	if dotSize <= 0 || len(arr) <= 1 {
		return arr[0]
	}

	if len(arr[1]) < dotSize {
		return fData
	}

	eStr := arr[1]
	endStr := eStr[:dotSize]
	ret := arr[0] + "." + endStr
	return ret
}

// 현재가격을 입력 => 수량 리턴
func (ty *ExchangeLotTickInfo) ConvertPriceToVol(posPrice string, convPrice float64) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error recover ", err)
		}
	}()

	if posPrice == "" || convPrice <= 0 {
		return "0"
	}

	fPosPrice, _ := strconv.ParseFloat(posPrice, 64)

	//사용가능 수량 계산
	fVol := convPrice / fPosPrice

	if ty.lotDotSize <= 0 {
		nVol := int64(math.Floor(fVol))
		return fmt.Sprintf("%d", nVol)
	}

	//-- 원래는 틱사이즈별로 계산해야하지만 => 바이넨스가 대부분 소숫점 단위로 계산한다 그래서 소숫점 아래 삭제하는방식으로
	newf := math.Floor(fVol/ty.fLotSize) * ty.fLotSize
	fmtstr := fmt.Sprintf("%%0.%df", ty.lotDotSize)
	strVol := fmt.Sprintf(fmtstr, newf)
	return strVol
}
