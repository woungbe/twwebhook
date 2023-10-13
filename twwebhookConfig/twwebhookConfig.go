package twwebhookConfig

import (
	"encoding/json"
	"fmt"
	"os"
)

var debugmodeflg bool = false //디버그 모드 플레그
type SvrConfigData struct {
	ServerName    string        //서버명
	RunMode       string        // = release (서비스용) , debug (디버그용)
	LogFile       string        //= Y/N  로그정보 파일로 저장할것인가
	Ssluse        string        //= Y,N
	Sslkey        string        //= ssl key 파일
	Sslcrt        string        //= ssl crt 파일
	RestLinkURL   string        // 서비스 루트 링크 ("" 인경우 디폴트로 "api")
	RestPort      int           //	HTTP 포트 번호
	CORSuse       bool          //크로스 도메인 요청 처리할것인가 ( 외부 html을 사용시 true 할것)
	LogSvrAddr    string        //로그서버 접속 주소 (127.0.01:000/Log  형태의 포트포함 URL까지)
	LogSvrUse     bool          //로그를 서버로 전송할것인가.
	LogSendLevel  int           //로그전송 레벨 ( 0= 전체, 1= 워링 부터)
	TelegramToken string        //텔레그램 bot 토큰
	LimitorUse    bool          //리미터 사용여부
	LimitorInfo   LimitorConfig //리미터 설정
}
type LimitorConfig struct {
	UrlPathLength        int64 //urlPath 체크 타임 ( 1초 단위로 설정,  )
	AllUIDLength         int64 //UID당 체크 타임 ( 1초 단위로 = 기본값은 1분= 60)
	UrlPathLimitOverTime int64 //urlpath 언락 타임 ( 1초 단위)
	AllUIDLimitOverTime  int64
	UrlPathLimitCount    int //urlpath 리밋 카운터
	AllUIDLimitCount     int
}

type APISvrConfig struct {
	cnf           SvrConfigData
	configOrgData string //설정 정보 데이터
}

var APISvrConfigPtr *APISvrConfig

// IsDebugmode 현재  디버그 모드이가 체크
func IsDebugmode() bool {
	return debugmodeflg
}

// SetDebugmode 디버그모드 설정
func SetDebugmode(b bool) {
	debugmodeflg = b
}

// IsCROS 외부 파일로 실행시 Y할것
func IsCROS() bool {
	return APISvrConfigPtr.cnf.CORSuse
}

func GetConfig() *APISvrConfig {

	if APISvrConfigPtr == nil {
		APISvrConfigPtr = new(APISvrConfig)
	}
	return APISvrConfigPtr
}

// GetServerID 서버 ID 가져오기
func GetServerID() string {
	return APISvrConfigPtr.cnf.ServerName
}

func StartLimitor() {
	APISvrConfigPtr.cnf.LimitorUse = true
}
func StopLimitor() {
	APISvrConfigPtr.cnf.LimitorUse = false
}

// GetConfigData 서버 기본정보
func GetConfigData() *SvrConfigData {
	return &APISvrConfigPtr.cnf
}

func (ty *APISvrConfig) InitConfig(isDebug bool) error {

	if isDebug {
		er := ty.loadConfig_Debug()

		if er != nil {
			return er
		}
	}
	// else {
	// 	er := ty.loadConfig()
	// 	if er != nil {
	// 		return er
	// 	}
	// }

	return nil
}

func (ty *APISvrConfig) GetConfigOrgData() string {
	return ty.configOrgData
}

// loadConfig_Debug 암호화 되지않은 설정파일
func (ty *APISvrConfig) loadConfig_Debug() error {

	fmt.Println("========================== Warning ==========================================")
	fmt.Println("Warn", "설정파일 암호화 안됌", "실서비스인경우 암호화된 설정파일을 사용하십시요!!!!")
	fmt.Println("========================== Warning ==========================================")
	b, err := os.ReadFile("./webhookconfig.json")
	if err != nil {
		fmt.Println("Warn", "config file Not found", "webhookconfig.json")
		return err
	}
	ty.configOrgData = string(b)

	er := json.Unmarshal(b, &ty.cnf)
	if er != nil {
		fmt.Println("Error", "설정로드에러", er.Error())
		return er
	}

	if ty.cnf.RunMode == "debug" {
		SetDebugmode(true)
	} else {
		SetDebugmode(false)
	}

	return nil
}

// loadConfig 암호화 되어있는 설정 파일
// func (ty *APISvrConfig) loadConfig() error {

// b, err := os.ReadFile("./webhookconfig.aes")
// if err != nil {

// 	fmt.Println("Warn", "config file Not found", "webhookconfig.json")
// 	return err
// }

// //-- 복호화
// dncData, err2 := dfinUtil.XFAESdec(dfinDefine.AESKey, b)
// if err2 != nil {
// 	pawfmt.Println("Error", "config file decode error", err2.Error())
// 	return err2
// }

// ty.configOrgData = string(dncData)

// er := json.Unmarshal([]byte(dncData), &ty.cnf)
// if er != nil {
// 	fmt.Println("Error", "설정로드에러", er.Error())
// 	return er
// }

// if ty.cnf.RunMode == "debug" {
// 	SetDebugmode(true)
// } else {
// 	SetDebugmode(false)
// }

// return nil
// }
