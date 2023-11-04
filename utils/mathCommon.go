package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// 소수점 몇째자리인가 카운트 하는것
func DecimalCount(numbers string) (bool, int) {
	value, err := strconv.ParseFloat(numbers, 64)
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

func CeilFloat(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	Ceil := math.Ceil(num*shift) / shift
	return Ceil
}

// 현재가, 최소금액
func MinCoinSize(curPrice, minMoney, minSize string, coinStep int) float64 {

	mCurPrice, err := strconv.ParseFloat(curPrice, 64)
	if err != nil {
		fmt.Println(err)
	}

	mMinMoney, err := strconv.ParseFloat(minMoney, 64)
	if err != nil {
		fmt.Println(err)
	}

	mMinSize, err := strconv.ParseFloat(minSize, 64)
	if err != nil {
		fmt.Println(err)
	}

	// mCoinStep, err := strconv.ParseFloat(coinStep, 64)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	if mMinMoney/mCurPrice < mMinSize {
		return mMinSize
	} else {
		x := mMinMoney / mCurPrice
		cnt := math.Pow(10, float64(coinStep))
		return math.Ceil(x*cnt) / cnt
	}
}

func Cell(num float64, decimalCnt string) float64 {
	tmp, er := strconv.Atoi(decimalCnt)
	if er != nil {
		fmt.Println("roundToTwoDecimalPlaces :", er)
	}

	cnt := math.Pow(10, float64(tmp))
	return math.Floor(num*cnt) / cnt
}

func roundToTwoDecimalPlaces(num float64, decimalCnt string) float64 {
	tmp, er := strconv.Atoi(decimalCnt)
	if er != nil {
		fmt.Println("roundToTwoDecimalPlaces :", er)
	}

	cnt := math.Pow(10, float64(tmp))
	return math.Floor(num*cnt) / cnt
}
