package main

func main() {

	init()

}

// 전체 초기확
func init() {

	// 마켓데이터 가져오기 - 보관용 , 자릿수, 최소최대 수량 구조
	GetMarketData()

	// 유저정보 DB로 가져오기 !! - 보관용, 주문 처리관련 데이터도 있음
	GetDBInit()

	// http - 터미널 조율을. 외부 조율로 변경해서 처리 요청
	http()

	// 로직
	ctrl.Init()

}

// 마켓 데이터 처리
func GetMarketData()

// DB초기화 - 마켓데이터 호출
func GetDBInit()

// http 초기화
func http()
