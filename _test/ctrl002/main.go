package main 

/*

해당 신호가 들어오면,  




*/

type Ctrl struct {
	StretegyList map[string][]UsersStretegy // 전략별 유저 전략측정 
}



// 데이터 받기 
func (ty *Ctrl) DataReceive(data interface{}){
	// 어떤 전략인지 검증해서 찾기
	userStretegy := StretegyList[data.type]
	userStretegy.Actions()
}

