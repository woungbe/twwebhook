package main

import (
	"fmt"
	"os"
	binanceorder "twwebhook/binanceOrder"
)

func main() {
	b := Client()
	// GetExchange(b)    // 마켓정보 가져오기
	// GetServiceTime(b) // 현재시간 가져오기
	GetLeverageList(b) // 레버리지 정보 가져오기
}

func Client() *binanceorder.BinanceInfo {

	err := godotenv.Load(".env")
	b := new(binanceorder.BinanceInfo)
	key1 := os.Getenv("public")
	key2 := os.Getenv("secrit")
	b.Init(key1, key2)
	return b
}

func GetExchange(b *binanceorder.BinanceInfo) {
	ex, er := b.GetExchangeInfo()
	if er != nil {
		fmt.Println(er)
	}

	for _, v := range ex.Symbols {
		fmt.Printf("%+v\n", v)
	}
}

func GetServiceTime(b *binanceorder.BinanceInfo) {
	n, er := b.GetServiceTime()
	if er != nil {
		fmt.Println(er)
	}

	fmt.Println("GetServiceTime : ", n)
}

func GetLeverageList(b *binanceorder.BinanceInfo) {
	leverageBracket, er := b.GetLeverageList()
	if er != nil {
		fmt.Println(er)
	}

	for _, v := range leverageBracket {
		fmt.Println("GetServiceTime : ", v)
	}

}
