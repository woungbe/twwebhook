package main

import (
	"fmt"
	"os"
)

func main() {

	err := godotenv.Load(".env")

	bexinfo := new()
	bexinfo.Init(
		os.Getenv("public"),
		os.Getenv("secrit"),
	)

	bexinfo.GetMarketInfo()
	bexinfo.GetLeverageList()

	fmt.Println("000000000000000000000000000000000")

	done := make(chan bool)
	<-done

}
