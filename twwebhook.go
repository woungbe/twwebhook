package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"twwebhook/twwebhookConfig"
	"twwebhook/twwebhookService"
)

var pgVersion = "Tradingview WebHook API V.1.0.0 - Release [2023-08-07]"

func main() {
	InitSrevice()

	done := make(chan bool)
	r := bufio.NewReader(os.Stdin)
	go func() {
		for {
			line, err := r.ReadString('\n')
			if err != nil && err.Error() != "unexpected newline" {
				fmt.Println(err.Error())
				//	return
				line = ""
			}

			line = strings.TrimSpace(line)
			CMDPaser(line)
		}
	}()
	<-done
}

// InitSrevice 서비스 초기화
func InitSrevice() {

	cnf := twwebhookConfig.GetConfig()
	er := cnf.InitConfig(true) //실서비스때는 꼭 false로 처리할것

	if er != nil {
		fmt.Println("설정정보로드에러")
		fmt.Println(er)
		os.Exit(1)
	}

	// if twwebhookConfig.GetConfigData().LogFile == "Y" {
	// 	// pawlog.CreateLogSystem(5, true, "logs/", true)
	// } else {
	// 	// pawlog.CreateLogSystem(5, false, "", false)
	// }

	//-- 에러 컨트롤러 초기화
	// var logSVRinfo // pawlog.MainFrameConnectInfo
	// logSVRinfo.ServerAddress = twwebhookConfig.GetConfigData().LogSvrAddr
	// logSVRinfo.Ssluse = "N" // svrconfig.GetConfigData().Ssluse
	// logSVRinfo.Sslkey = ""  //svrconfig.GetConfigData().Sslkey
	// logSVRinfo.Sslcrt = ""  //svrconfig.GetConfigData().Sslcrt

	// if twwebhookConfig.GetConfigData().LogSvrUse {
	// 	// er2 := pawlog.GetXFError().Init(twwebhookConfig.GetConfigData().ServerName, logSVRinfo, true)
	// 	if er2 != nil {
	// 		log.Error("Error", "Covest Pro Log Server 초기화 오류", er2.Error())
	// 		//os.Exit(1)
	// 	} else {
	// 		// pawlog.SetNetSendType(twwebhookConfig.GetConfigData().LogSendLevel)
	// 		// pawlog.GetXFError().StartLogSending()
	// 	}
	// } else {
	// 	er2 := // pawlog.GetXFError().Init(twwebhookConfig.GetConfigData().ServerName, logSVRinfo, false)
	// 	if er2 != nil {
	// 		log.Error("Error", "Covest Pro Log Server 초기화 오류", er2.Error())
	// 		//os.Exit(1)
	// 	}
	// }

	if twwebhookConfig.IsDebugmode() {
		fmt.Println("================= DEBUG MODE !! =============================")
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	twwebhookService.GetTelegram().StartTelegram()

	//REST API Service 시작

	twwebhookService.HttpInit()

	fmt.Println("================= Server Start OK  =============================")
	fmt.Println("=================  Server Version  =============================")
	fmt.Println(pgVersion)
	fmt.Println("================================================================")
}

// CMDPaser 콘솔 명령 파서
func CMDPaser(strCMD string) {
	defer func() {
		if err := recover(); err != nil {
			// pawlog.Error("Crit Panic", "Error", err)
		}
	}()
	if strCMD == "" {
		return
	}
	if strCMD == "exit" {
		os.Exit(1)
		return
	}
}
