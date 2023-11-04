package twebdefined

/*
유저 등록
유저 조회
유저 삭제 -- DB에서 해

전략 등록
전략 조회
전략 삭제

유저 - 전략 매칭 등록
유저 - 전략 매칭 수정
유저 - 전략 조회
유저 - 전략 정지
유저 - 전략 시작


바이낸스 정보 강제 갱신
: 마켓데이터 재호출
: 바스켓 데이터 재호출

// 유저 전략 조회 - 시작도 있음

*/

type REST_TW_Default struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 유저추가 등록 reqeust
type REST_TW_AddUser struct {
	Name   string `json:"name"`   // 이름
	Access string `json:"access"` // public key
	Script string `json:"script"` // script key
}

// 유저추가 등록 response
type REST_TW_AddUserRES struct {
	Useridx int `json:"useridx"` // 유저고유번호
}

// 유저 검색 reqeust
type REST_TW_GetUser struct {
	Name string `json:"name"` // 이름
}

// 유저 검색 response
type REST_TW_GetUserRES struct {
	Useridx int    `json:"useridx"`
	Name    string `json:"name"`   // 이름
	Access  string `json:"access"` // public key
	Script  string `json:"script"` // script key
}

// 전략 등록 request
type REST_TW_AddStrategy struct {
	Stname    string `json:"stname"`
	Stcontent string `json:"stcontent"`
}

// 전략 등록 response
type REST_TW_AddStrategyRES struct {
	REST_TW_Default
}

// 전략 검색 strategy reqeust
type REST_TW_Strategy struct {
	SearchText string `json:"searchText"` // 검색어, 없으면 전체
}

// 전략 검색 strategy response
type REST_TW_StrategyRES struct {
	Strategysrl int    `json:"strategysrl"` // 고유번호
	StName      string `json:"stName"`      // 이름
	StContent   string `json:"stContent"`   // 설명
}

// 전략 삭제 reqeust
type REST_TW_DelStrategy struct {
	Strategysrl int `json:"strategysrl"` // 고유번호
}

// 전략 삭제 response
type REST_TW_DelStrategyRES struct {
	REST_TW_Default
}

// 전략 매칭 등록
// order setting reqeust
type REST_TW_AddMapping struct {
	Useridx     int    `json:"useridx"`     // 유저고유번호
	Strategysrl int    `json:"strategysrl"` // 전략 고유번호
	Profitflg   int    `json:"profitflg"`   // 익절 여부
	Profitval   string `json:"profitval"`   // 익절 값
	Losscutflg  int    `json:"losscutflg"`  // 로스컷 여부
	Losscutval  string `json:"losscutval"`  // 로스컷 값
	Rateflg     int    `json:"rateflg"`     // 비율/고정 설정
	Rateval     int    `json:"rateval"`     // 비율 값
	Fixedval    int    `json:"fixedval"`    // 고정 값
	Counterflg  int    `json:"counterflg"`  // 0:아무것도 안함 , 1: 롱,숏일때 청산,매수
	Startflg    int    `json:"Startflg"`    // 0:정지, 1: 시작
}

// 맵핑 설정 리턴
// order setting response
type REST_TW_AddMappingRES struct {
	REST_TW_Default
}

// 전략 매칭 수정 reqeust
type REST_TW_UpdateMapping struct {
	Mappingsrl  int    `json:"mappingsrl"`  // 매칭 고유번호
	Useridx     int    `json:"useridx"`     // 유저고유번호
	Strategysrl int    `json:"strategysrl"` // 전략 고유번호
	Profitflg   int    `json:"profitflg"`   // 익절 여부
	Profitval   string `json:"profitval"`   // 익절 값
	Losscutflg  int    `json:"losscutflg"`  // 로스컷 여부
	Losscutval  string `json:"losscutval"`  // 로스컷 값
	Rateflg     int    `json:"rateflg"`     // 비율/고정 설정
	Rateval     int    `json:"rateval"`     // 비율 값
	Fixedval    int    `json:"fixedval"`    // 고정 값
	Counterflg  int    `json:"counterflg"`  // 0:아무것도 안함 , 1: 롱,숏일때 청산,매수
	Startflg    int    `json:"Startflg"`    // 0:정지, 1: 시작
}

// 전략 매칭 수정 response
type REST_TW_SetMappingRES struct {
	REST_TW_Default
	Mappingsrl int `json:"mappingsrl"`
}

// 유저 전략 조회 reqeust
type REST_TW_Mapping struct {
	Useridx int `json:"useridx"`
}

// 유저 전략 조회 response
type REST_TW_MappingRES struct {
	Mappingsrl  int    `json:"mappingsrl"`  // 매칭 고유번호
	Useridx     int    `json:"useridx"`     // 유저고유번호
	Strategysrl int    `json:"strategysrl"` // 전략 고유번호
	Profitflg   int    `json:"profitflg"`   // 익절 여부
	Profitval   string `json:"profitval"`   // 익절 값
	Losscutflg  int    `json:"losscutflg"`  // 로스컷 여부
	Losscutval  string `json:"losscutval"`  // 로스컷 값
	Rateflg     int    `json:"rateflg"`     // 비율/고정 설정
	Rateval     int    `json:"rateval"`     // 비율 값
	Fixedval    int    `json:"fixedval"`    // 고정 값
	Counterflg  int    `json:"counterflg"`  // 0:아무것도 안함 , 1: 롱,숏일때 청산,매수
	Startflg    int    `json:"Startflg"`    // 0:정지, 1: 시작
}

// 유저 전략 정지 reqeust
type REST_TW_StopStrategy struct {
	Mappingsrl int `json:"mappingsrl"` // 매칭 고유번호
	Useridx    int `json:"useridx"`    // 유저고유번호
}

// 유저 전략 정지 response
type REST_TW_StopStrategyRES struct {
	REST_TW_Default
}

// 유저 전략 시작 reqeust
type REST_TW_StartStrategy struct {
	Mappingsrl int `json:"mappingsrl"` // 매칭 고유번호
	Useridx    int `json:"useridx"`    // 유저고유번호
}

// 유저 전략 시작 response
type REST_TW_StartStrategyRES struct {
	REST_TW_Default
}
