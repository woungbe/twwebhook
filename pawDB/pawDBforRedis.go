package pawDB

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisDbConnectInfo 데이터베이스 연결 정보
type RedisDbConnectInfo struct {
	StrIP     string
	StrPasswd string
	StrDBname int
}

type RedisDB struct {
	stDBConnectInfo RedisDbConnectInfo
	redisDB         *redis.Client
}

// 이걸 클래스로 어떻게 만드냐 ?..
func (ty *RedisDB) Init(tyDBinfoStruct RedisDbConnectInfo) error {

	ty.stDBConnectInfo = tyDBinfoStruct
	ty.redisDB = redis.NewClient(&redis.Options{
		Addr:         ty.stDBConnectInfo.StrIP,
		Password:     ty.stDBConnectInfo.StrPasswd,
		DB:           ty.stDBConnectInfo.StrDBname,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolSize:     10,
		PoolTimeout:  10 * time.Second,
	})

	if ty.redisDB == nil {
		return fmt.Errorf("redis init error")
	}

	return nil
}

func (ty *RedisDB) IsExistsKey(key string) (bool, error) {

	stringCmd := ty.redisDB.Exists(key)
	n, err := stringCmd.Result()
	var retB bool
	retB = false
	if n > 0 {
		retB = true
	}
	return retB, err
}

func (ty *RedisDB) SetString(key string, value string, TTL int64) error {
	delTime := time.Duration(TTL) * time.Second
	err := ty.redisDB.Set(key, value, delTime).Err()
	if err != nil {
		//panic(err)
		return err
	}

	return nil
}

func (ty *RedisDB) SetRightPushList(key string, value string) (int64, error) {
	IntCmd := ty.redisDB.RPush(key, value)
	cnt, err := IntCmd.Result()

	if err != nil {
		return 0, err
	}

	return cnt, err
}

func (ty *RedisDB) SetLfitPushList(key string, value string) (int64, error) {
	IntCmd := ty.redisDB.LPush(key, value)
	cnt, err := IntCmd.Result()

	if err != nil {
		return 0, err
	}

	return cnt, err
}

// source, destination string, timeout time.Duration
func (ty *RedisDB) LTrim(key string, start, stop int64) (string, error) {
	StatusCmd := ty.redisDB.LTrim(key, start, stop)
	str, err := StatusCmd.Result()

	if err != nil {
		return "", err
	}

	return str, err
}

func (ty *RedisDB) Setset(key string, value ...string) (int64, error) {
	// value가 같다면 중복저장안됨
	// 저장이 안됐다면 checkData == 0 ,
	// 저장됐다면. checkData == 1 이상,
	IntCmd := ty.redisDB.SAdd(key, value)
	//str := IntCmd.String()
	checkData, err := IntCmd.Result()

	if err != nil {
		return 0, err
	}

	if checkData == 0 {
		return 0, fmt.Errorf("저장되지 않았습니다. ")
	}

	return checkData, nil
	//setDo("sadd", key, value)

}

// 필드가 중복되면 들어가지않음,
// 저장 성공시 b = true
// 저장 실패시 b = false
func (ty *RedisDB) SetHash(key string, field string, value interface{}) (bool, error) {
	bolCMD := ty.redisDB.HSetNX(key, field, value)
	b, err := bolCMD.Result()
	return b, err
}

// value 가 같으면 안들어감, cnt = 0, score는 상관없음
func (ty *RedisDB) SetZset(key string, score float64, value interface{}) (int64, error) {
	// Zcmd := &redis.Z{Member: value, Score: score}
	// zSet := ty.redisDB.ZAdd(key, &Zcmd)
	// cnt, err := zSet.Result()
	// if err != nil {
	// 	return cnt, err
	// }
	return 0, nil
}

// 저장하는 방법이나 사용법을 모름
// func SetStrem() error { return nil }
func (ty *RedisDB) setDo(data ...string) error {
	_, err := ty.redisDB.Do(data).String()
	if err != nil {
		return err
	}
	return nil
}

func (ty *RedisDB) GetString(key string) (string, error) {

	stringCmd := ty.redisDB.Get(key)
	n, err := stringCmd.Result()

	if err != nil {
		return "", err
	}

	return n, err
}

func (ty *RedisDB) GetStringKeyList(key string, cnt int64) []string {

	var cursor uint64
	var n int
	var send []string

	for {
		var keys []string
		var err error
		keys, cursor, err = ty.redisDB.Scan(cursor, key, cnt).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		send = append([]string{}, append(send, keys...)...)

		if cursor == 0 {
			break
		}
	}
	return send
}

func (ty *RedisDB) DelString(key string) (int64, error) {

	IntCmd := ty.redisDB.Del(key)
	n, err := IntCmd.Result()

	if err != nil {
		return 0, err
	}

	return n, err
}

// 리스트 갯수 찾기
func (ty *RedisDB) GetListLen(key string) (int64, error) {

	IntCmd := ty.redisDB.LLen(key)
	n, err := IntCmd.Result()

	if err != nil {
		return 0, err
	}

	return n, err
}

// 리스트 값만 가져오기
func (ty *RedisDB) GetListVal(key string, start, stop int64) ([]string, error) {

	startSlideCmd := ty.redisDB.LRange(key, start, stop)
	str, err := startSlideCmd.Result()

	if err != nil {
		return nil, err
	}

	return str, err
}

func (ty *RedisDB) GetListKeyList(key string, cnt int64) []string {
	return ty.GetStringKeyList(key, cnt)
}

func (ty *RedisDB) DelList(key string) (int64, error) {
	return ty.DelString(key)
}

// hash중에 하나의 필드만 가져오는 것
func (ty *RedisDB) GetHash(key string, field string) (string, error) {

	StringSliceCmd := ty.redisDB.HGet(key, field)
	str, err := StringSliceCmd.Result()
	if err != nil {
		return "", err
	}

	return str, err
}

// key 리스트만 가져오는 것
func (ty *RedisDB) GetHashKeyList(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	//ScanCmd := rsdb.HScan(key, cursor, "", count)

	ScanCmd := ty.redisDB.HScan(key, cursor, match, count)

	str, cursor, err := ScanCmd.Result()
	if err != nil {
		return nil, 0, err
	}
	return str, cursor, err
}

// value만 가져오는 것
func (ty *RedisDB) GetHashValueList(key string) ([]string, error) {
	StringSliceCmd := ty.redisDB.HVals(key)
	str, err := StringSliceCmd.Result()
	if err != nil {
		return nil, err
	}

	return str, err
}

// key value로 map형태로 조합해서 만들기
func (ty *RedisDB) GetHashKeyValueMap(key string, ncursor int64, match string, count int64) (map[string]string, int64, error) {
	cursor := uint64(ncursor)

	StrCmd := ty.redisDB.HScan(key, cursor, match, count)
	// keys []string, cursor uint64, err error
	str, next, err := StrCmd.Result()
	if err != nil {
		return nil, 0, err
	}
	HashData := make(map[string]string)

	for i := 0; i < len(str); i++ {
		//HashData.put(str[i], str[i+1])
		HashData[str[i]] = str[i+1]
		i++
	}

	nextcnt := int64(next)
	return HashData, nextcnt, err
}

// 필드 하나 삭제
func (ty *RedisDB) DelHashField(key string, field string) (int64, error) {
	strCmd := ty.redisDB.HDel(key, field)
	cnt, err := strCmd.Result()
	if err != nil {
		return 0, err
	}

	return cnt, err
}

// 필드 전체 삭제
func (ty *RedisDB) DelHashKey(key string) (int64, error) {
	return ty.DelString(key)
}

// 필드 하나를 업데이트
// true => 추가된거
// false , nil => 업데이트 된거
// false , !nil => 문제가 있는거
func (ty *RedisDB) UpdateHashField(key string, field string, val interface{}) (bool, error) {

	// key, field string, value interface{}
	strCmd := ty.redisDB.HSet(key, field, val)
	b, err := strCmd.Result()
	if err != nil {
		return false, err
	}

	return b, err
}

// zset 전체 rows 가져오기
func (ty *RedisDB) GetZsetLen(key string) (int64, error) {

	intCmd := ty.redisDB.ZCard(key)
	cnt, err := intCmd.Result()
	if err != nil {
		return 0, err
	}

	return cnt, err
}

// value로 스코어 출력하기
func (ty *RedisDB) GetZsetValToScore(key, value string) (float64, error) {

	floatCmd := ty.redisDB.ZScore(key, value)
	cnt, err := floatCmd.Result()
	if err != nil {
		return 0, err
	}

	return cnt, err

}

// value 만 출력하기
func (ty *RedisDB) GetZsetTotalVal(key string, start, end int64) ([]string, error) {

	StringSliceCmd := ty.redisDB.ZRange(key, start, end)
	str, err := StringSliceCmd.Result()
	if err != nil {
		return nil, err
	}

	return str, err
}

// value score 둘다 출력하기
func (ty *RedisDB) GetZsetAll(key string, start, end int64) (map[string]float64, error) {
	//zrange key:zset 0 3 withscores

	ZSliceCmd := ty.redisDB.ZRangeWithScores(key, start, end)
	str, err := ZSliceCmd.Result()
	if err != nil {
		return nil, err
	}

	data := make(map[string]float64)

	for _, v := range str {
		val := v.Member.(string)
		data[val] = v.Score
	}

	return data, err

}

// zset value로 삭제하기
func (ty *RedisDB) DelZsetVal(key, field string) (int64, error) {

	IntCmd := ty.redisDB.ZRem(key, field)
	cnt, err := IntCmd.Result()
	if err != nil {
		return 0, err
	}

	return cnt, err
}

// zset score로 지우기
func (ty *RedisDB) ZRemRangeByScore(key, min, max string) (int64, error) {

	IntCmd := ty.redisDB.ZRemRangeByScore(key, min, max)
	cnt, err := IntCmd.Result()
	if err != nil {
		return 0, err
	}

	return cnt, err
}

func (ty *RedisDB) DelZsetAll(key string) (int64, error) {
	return ty.DelString(key)
}

func (ty *RedisDB) GetSetAll(key string, ncursor int64, match string, count int64) ([]string, int64, error) {

	cursor := uint64(ncursor)
	ScanCmd := ty.redisDB.SScan(key, cursor, match, count)
	str, nextcs, err := ScanCmd.Result()
	if err != nil {
		return nil, 0, err
	}

	send := int64(nextcs)
	return str, send, err

}

func (ty *RedisDB) DelSetVal(key string, val string) (int64, error) {

	ScanCmd := ty.redisDB.SRem(key, val)
	cnt, err := ScanCmd.Result()
	if err != nil {
		return 0, err
	}

	return cnt, err

}

func (ty *RedisDB) DelSetKey(key string) (int64, error) {
	return ty.DelString(key)
}
