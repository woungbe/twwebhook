package twwebhookService

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"twwebhook/twwebhookConfig"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type JSONuserinfo struct {
	UserList []int64
}
type TelegramCtrl struct {
	pBot    *tgbotapi.BotAPI
	bInited bool //초기화 상태
	//UserIDX int64 //사용자 dix
	userIdxLst map[int64]int64
}

var TelegramCtrlPTR *TelegramCtrl

// GetTelegram 텔레그램 인스턴스
func GetTelegram() *TelegramCtrl {
	if TelegramCtrlPTR == nil {
		TelegramCtrlPTR = new(TelegramCtrl)
		TelegramCtrlPTR.Init()
	}

	return TelegramCtrlPTR
}

// Init 텔레그램 초기화
func (ty *TelegramCtrl) Init() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Crit Panic", "Error", err)
		}
	}()
	var err error

	ty.userIdxLst = make(map[int64]int64)

	er := ty.loadUserList()
	if er != nil {
		fmt.Println("Error", "텔레그램 사용자IDX 로드오류", er.Error())
	}

	cnf := twwebhookConfig.GetConfigData()
	ty.pBot, err = tgbotapi.NewBotAPI(cnf.TelegramToken)
	if err != nil {
		fmt.Println("Error", "텔레그램 초기화오류", err.Error())
		return err
	}
	ty.bInited = true //초기화 상태

	//go ty.procFunc()
	return nil
}

func (ty *TelegramCtrl) StartTelegram() {
	if ty.bInited == false {
		return
	}
	go ty.procFunc()
}
func (ty *TelegramCtrl) procFunc() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Crit Panic", "Error", err)
		}
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := ty.pBot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		receiveMsg := update.Message.Text
		if len(receiveMsg) > 3 && receiveMsg[0] == '/' {
			// 명령부와 데이터 부 나눔
			userIdx := update.Message.Chat.ID
			cmds := strings.Split(receiveMsg, " ")

			if len(cmds) > 0 {
				if cmds[0] == "/reg" || cmds[0] == "/등록" {

					ty.userIdxLst[userIdx] = userIdx
					ty.saveUserList()
					ty.SendMsg("등록되었습니다.", userIdx)
				}
			}
		}
		runtime.Gosched()
	}
}

func (ty *TelegramCtrl) SendMsgAll(msg string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Crit Panic", "Error", err)
		}
	}()
	for _, value := range ty.userIdxLst {
		ty.SendMsg(msg, value)
	}
	return nil
}

// SendMsg 메시지 샌드
func (ty *TelegramCtrl) SendMsg(msg string, userIdx int64) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Crit Panic", "Error", err)
		}
	}()

	botmsg := tgbotapi.NewMessage(userIdx, msg)
	_, er := ty.pBot.Send(botmsg)
	if er != nil {
		strEr := er.Error()
		if "Forbidden: bot was blocked by the user" == strEr {
			//- 사용자가 임의로 삭제한경우
			//fmt.Println(strEr)
		}
	}
	return er
}

// saveUserList 등록된 사용자 정보 저장
func (ty *TelegramCtrl) saveUserList() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Crit Panic", "Error", err)
		}
	}()
	var ulst JSONuserinfo
	for _, value := range ty.userIdxLst {
		ulst.UserList = append(ulst.UserList, value)
	}

	resultData, _ := json.Marshal(ulst)
	err := os.WriteFile("./telegramUserList.json", resultData, os.FileMode(644))
	if err != nil {
		fmt.Println("텔레그램 유저 저장 실패")
		return err
	}

	return nil
}

func (ty *TelegramCtrl) loadUserList() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Crit Panic", "Error", err)
		}
	}()
	var ulst JSONuserinfo
	b, err := os.ReadFile("./telegramUserList.json")
	if err != nil {
		return err
	}

	json.Unmarshal(b, &ulst)
	for _, v := range ulst.UserList {
		ty.userIdxLst[v] = v
	}
	return nil
}
