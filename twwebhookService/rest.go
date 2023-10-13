package twwebhookService

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"twwebhook/twwebhookConfig"
)

// 서비스용 HTTP
type APIHttp struct {
	strLink        string
	strPort        string
	mMux           *http.ServeMux
	mServiceStatus int //서비스 상태  1= 정상, 100= 서비스 정지상태

	CheckLastConErr int
}

var HttpSvr *APIHttp
var gLockUser map[string]string //리미터에 걸러진 유저중 블럭처리된 유저맵

func GetHttpSvr() *APIHttp {
	return HttpSvr
}

func HttpInit() {
	if HttpSvr == nil {
		HttpSvr = new(APIHttp)
		HttpSvr.init()
	}

	gLockUser = make(map[string]string)

	cnf := twwebhookConfig.GetConfigData()
	// respLimitor.GetLimitorCounterCtrl().Init()
	// respLimitor.GetLimitorCounterCtrl().SetConfig(time.Second*time.Duration(cnf.LimitorInfo.UrlPathLength), time.Second*time.Duration(cnf.LimitorInfo.AllUIDLength), time.Second*time.Duration(cnf.LimitorInfo.UrlPathLimitOverTime), time.Second*time.Duration(cnf.LimitorInfo.AllUIDLimitOverTime), cnf.LimitorInfo.UrlPathLimitCount, cnf.LimitorInfo.AllUIDLimitCount)

	if cnf.RestLinkURL == "" {
		HttpSvr.strLink = "api"
	} else {
		HttpSvr.strLink = cnf.RestLinkURL
	}

	if cnf.RestPort <= 0 {
		HttpSvr.strPort = "80"
	} else {
		HttpSvr.strPort = strconv.Itoa(cnf.RestPort)
	}
	HttpSvr.mServiceStatus = 1

	go HttpSvr.HTTPStart()
}

// 서비스 정지상태로
func StopServiceStatus() {
	HttpSvr.mServiceStatus = 100
}
func StartServiceStatus() {
	HttpSvr.mServiceStatus = 1
}

func IsServiceExit() bool {
	if HttpSvr.mServiceStatus == 1 {
		return false
	}
	return true
}

func getContentType(localPath string) string {
	var contentType string
	ext := filepath.Ext(localPath)

	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain"
	}

	return contentType
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
	defer func() {
		// 경과 시간
		if !twwebhookConfig.IsDebugmode() {
			if err := recover(); err != nil {
				fmt.Println("Crit Error !!!!! enableCors", "HTTP Error", err)
			}
		}
	}()

	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, PUT, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization , Access-Control-Allow-Headers , Origin , X-MBX-APIKEY ")
	(*w).Header().Add("Access-Control-Allow-Headers", "RecaptchaToken , Auth-Token , NewAuth-Token , TokenSHA1 , DeviceID , ClientTime , SeqData ")
	(*w).Header().Add("Access-Control-Allow-Headers", "X-RateLimit-Limit , Retry-After ")
	(*w).Header().Add("Access-Control-Expose-Headers", "NewAuth-Token")

}

//===============================================

func (ty *APIHttp) init() {
}
func (ty *APIHttp) GetMux() *http.ServeMux {
	return ty.mMux
}

// HTTPStart HTTP서비스 시작
func (ty *APIHttp) HTTPStart() {
	conf := twwebhookConfig.GetConfigData()
	ty.mMux = http.NewServeMux()

	ty.initRESTFunc()

	fmt.Println("Starting HTTP...")

	if conf.Ssluse == "Y" {
		fmt.Println("Starting server for HTTPS...")
		err := http.ListenAndServeTLS(":"+ty.strPort, conf.Sslcrt, conf.Sslkey, ty.mMux)
		if err != nil {
			fmt.Println("HTTP Error", "Http", err.Error())
			os.Exit(1)
			return
		}
	} else {
		fmt.Println("Starting server for HTTP...")
		err :=
			http.ListenAndServe(":"+ty.strPort, ty.mMux)
		if err != nil {
			fmt.Println("HTTP Error", "Http", err.Error())
			os.Exit(1)
			return
		}
	}
}

// rest api 함수 리스트 초기화
func (ty *APIHttp) initRESTFunc() {
	ty.mMux.HandleFunc("/twwebhook", RESTFunc_tradingviewwebhook) //서버 시간 갭 체크
}

// Handler 디폴트 핸들러
func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// 경과 시간
		if !twwebhookConfig.IsDebugmode() {
			if err := recover(); err != nil {
				fmt.Println("Crit Error !!!!! Panic Handler", "HTTP Error", err)
			}
		}
	}()

	if twwebhookConfig.IsCROS() {
		enableCors(&w, r)
	}

	var f string
	if r.URL.Path == "/" {
		f = "/index.html"
	} else {
		f = r.URL.Path
	}

	localPath := "/" + f //twwebhookConfig.GetHtmlRoot() + f
	content, err := os.ReadFile(localPath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}

	defer r.Body.Close()

	contentType := getContentType(localPath)
	w.Header().Add("Content-Type", contentType)

	w.Write(content)
}

//============

func IsMobile(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	return Is_mobile(ua)
}
func Is_mobile(useragent string) bool {
	// the list below is taken from
	// https://github.com/bcit-ci/CodeIgniter/blob/develop/system/libraries/User_agent.php
	defer func() {
		// 경과 시간
		if !twwebhookConfig.IsDebugmode() {
			if err := recover(); err != nil {
				fmt.Println("Crit Error !!!!! Panic Handler", "HTTP Error", err)
			}
		}
	}()

	mobiles := []string{
		"Mobile Explorer", "Palm", "Motorola", "Nokia", "Palm", "iPhone", "Apple iPhone", "iPad", "Apple iPod Touch", "Sony Ericsson", "Sony Ericsson", "BlackBerry", "O2 Cocoon", "Treo", "LG", "Amoi", "XDA", "MDA", "Vario", "HTC", "Samsung",
		"Sharp", "Siemens", "Alcatel", "BenQ", "HP iPaq", "Motorola", "PlayStation Portable", "PlayStation 3", "PlayStation Vita", "Danger Hiptop", "NEC", "Panasonic", "Philips", "Sagem", "Sanyo", "SPV", "ZTE", "Sendo", "Nintendo DSi", "Nintendo DS", "Nintendo 3DS", "Nintendo Wii", "Open Web", "OpenWeb", "Android", "Symbian", "SymbianOS", "Palm", "Symbian S60", "Windows CE", "Obigo", "Netfront Browser", "Openwave Browser", "Mobile Explorer", "Opera Mini", "Opera Mobile", "Firefox Mobile", "Digital Paths", "AvantGo", "Xiino", "Novarra Transcoder", "Vodafone", "NTT DoCoMo", "O2", "mobile", "wireless", "j2me", "midp", "cldc", "up.link", "up.browser", "smartphone", "cellphone", "Generic Mobile"}

	for _, device := range mobiles {
		if strings.Index(useragent, device) > -1 {
			return true
		}
	}
	return false
}

// func Limitor_HandleFunc(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if !twwebhookConfig.IsDebugmode() {
// 				if err := recover(); err != nil {
// 					fmt.Println("Crit Error !!!!! Panic ", "Error", err)
// 				}
// 			}
// 		}()

// 		//fmt.Println("==================", r.URL.Path)
// 		enableCors(&w, r)
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
// 		cnf := twwebhookConfig.GetConfigData()

// 		strIP, _, err := net.SplitHostPort(r.RemoteAddr)
// 		if err != nil {
// 			strIP = r.RemoteAddr
// 		}
// 		strIP = dfinUtil.CanonicalizeIP(strIP)
// 		realip := GetRealIP(r)
// 		urlpath := r.URL.Path
// 		if cnf.LimitorUse {

// 			ipKey := strIP + "/" + realip

// 			limitflg, utlcount, alluidCNT, uidf, _ := respLimitor.GetLimitorCounterCtrl().CheckLimit(ipKey, urlpath)
// 			//fmt.Println("카운터", ipKey, urlpath, limitflg, utlcount)
// 			if uidf {
// 				//UID로 제한 될경우 ==> 로그 저장 하자
// 				//msg := fmt.Sprintf("Real IP : %s , IP : %s , UserIDX=%d", realip, strIP, uidx)
// 				msg := fmt.Sprintf("Real IP : %s , IP : %s ", realip, strIP)
// 				fmt.Println("Warn", "RESP Limit !", msg)
// 			}
// 			if limitflg {
// 				if alluidCNT > (cnf.LimitorInfo.AllUIDLimitCount * 5) {
// 					//---> 이건 아예 차단처리 한다.
// 					_, s := gLockUser[ipKey]
// 					if !s {
// 						gLockUser[ipKey] = ipKey
// 						//msg := fmt.Sprintf("Real IP : %s , IP : %s , UserIDX=%d", realip, strIP, uidx)
// 						msg := fmt.Sprintf("Real IP : %s , IP : %s ", realip, strIP)
// 						fmt.Println("Warn", "RESP Limit! , Check the lock", msg)
// 					}
// 				}

// 				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", utlcount))
// 				//w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", 0))
// 				//w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", currentWindow.Add(l.windowLength).Unix()))
// 				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(time.Second*time.Duration(cnf.LimitorInfo.UrlPathLimitOverTime)))) // RFC 6585
// 				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
// 				return
// 			}
// 		}
// 		next(w, r)
// 	})
// }

// func GetRealIP(r *http.Request) string {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("Crit Error [ respLimitor.GetRealIP ]")
// 		}
// 	}()

// 	var ip string
// 	if tcip := r.Header.Get("True-Client-IP"); tcip != "" {
// 		ip = tcip
// 	} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
// 		ip = xrip
// 	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
// 		i := strings.Index(xff, ", ")
// 		if i == -1 {
// 			i = len(xff)
// 		}
// 		ip = xff[:i]
// 	} else {
// 		var err error
// 		ip, _, err = net.SplitHostPort(r.RemoteAddr)
// 		if err != nil {
// 			ip = r.RemoteAddr
// 		}
// 	}
// 	return dfinUtil.CanonicalizeIP(ip)
// }

func RESTFunc_tradingviewwebhook(w http.ResponseWriter, r *http.Request) {

	defer func() {
		r.Body.Close()
		if err := recover(); err != nil {
			fmt.Println("Crit Error ", "REST Error", err)
		}
	}()

	if IsServiceExit() {
		http.Error(w, "503 Service Stop", http.StatusServiceUnavailable)
		return
	}

	if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		strMsg := string(body)
		fmt.Println("트레이딩뷰 메시지", strMsg)
		GetTelegram().SendMsgAll(strMsg)
	}
	w.WriteHeader(http.StatusOK)
}

func RESTFunc_tradingviewwebhookJSON(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
		if !twwebhookConfig.IsDebugmode() {
			if err := recover(); err != nil {
				fmt.Println("Crit Error ", "REST Error", err)
			}
		}
	}()

	if IsServiceExit() {
		http.Error(w, "503 Service Stop", http.StatusServiceUnavailable)
		return
	}

	if r.Method == "POST" {
		var inf REST_WEB_TradingViewJson
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &inf) // .NewDecoder(r.Body).Decode(&inf)
		if err != nil {
			errmsg := fmt.Sprintf("형식 에러: %s", err.Error())
			GetTelegram().SendMsgAll(errmsg)
			return
		}

		send := inf.MakeText()
		fmt.Println(send)
		GetTelegram().SendMsgAll(send)
	}
	w.WriteHeader(http.StatusOK)
}

/*
{
	"strategyName":"슈퍼트레이더 트레이딩 시그널",
	"coin":"BTC",
	"positionSide":"LONG",
	"powerText":"(강력)",
	"leverage":"10",
	"marginType":"격리",
	"price":"21323.1",
	"profit":"1.03",
	"lostcut":"0.97",
	"bottomText":" ★항상, 매번 말씀드리지만, 꼭 손절을 지키면서 매매 진행하시길 바랍니다.
	( 본인 투자에 대한 책임은 본인에게 있습니다)
	💡누적 수익으로 항상 접근하셔서 뇌동매매를 방지하시길 바랍니다!"
}
*/
