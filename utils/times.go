package utils

import "time"

//GetCurrentTimestamp UTC 현재 시간 리턴
func GetCurrentTimestamp() int64 {
	return time.Now().UTC().Unix()
}
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
func GetCurrentDateToString(UTCflg bool) string {
	if UTCflg {
		now := time.Now().UTC()
		return now.Format("2006-01-02")
	}
	now := time.Now()
	return now.Format("2006-01-02")
}
