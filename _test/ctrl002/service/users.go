




var UserListPTR *UserList

type UserList {
	Acc map[int]UserInfo 
}

type Client interface{}

type UserInfo struct {
	useridx int
	key string 
	sec  string 
	client Client
}


func GetUserInfo() *UserList {
	UserListPTR
}

func (ty *UserList) init(){
	// 유저 DB에서 가져오고 
	// client 연결 처리까지 끝내고 저장하기 
}

func (ty *UserList) findClient(i int ) Client {
	// ty.Acc[i]
	// 있으면 리턴 해주기 
	return nil
}





key 
sec 
client 


