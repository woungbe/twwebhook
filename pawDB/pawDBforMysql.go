package pawDB

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
)

//DbConnectInfo 데이터베이스 연결 정보
type DbConnectInfo struct {
	StrID     string
	StrPasswd string
	StrIP     string
	NPort     int
	StrDBname string
	Readonly  bool //== selet만 가능한것인가.
}

type MetalScanner struct {
	valid bool
	value interface{}
}

type BCEDBforMysql struct {
	bInited         bool
	stDBConnectInfo DbConnectInfo //연결정보
	mDBconn         *sql.DB
	mMaxOpenConn    int
	mMaxIdelConn    int

	mJobTime time.Time
	idleTime int64
	closeFlg bool //데이터베이스 close됐는가
}

//InitDB 데이터베이스 초기화 한다. nMaxOpenConn =0 , nMaxIdelConn =0 으로 할경우 동적으로 설정된다.  , nIdleChkTime= 아이들 체크타임( 초단위)
func (ty *BCEDBforMysql) InitDB(tyDbInfoStruct DbConnectInfo, nMaxOpenConn int, nMaxIdelConn int, nIdleChkTime int64) {
	ty.stDBConnectInfo = tyDbInfoStruct
	ty.mMaxOpenConn = nMaxOpenConn
	ty.mMaxIdelConn = nMaxIdelConn

	ty.idleTime = nIdleChkTime
	ty.closeFlg = false
	if ty.idleTime < 60 {
		ty.idleTime = 60
	}

}

//ConnectDB 데이터베이스연결
func (ty *BCEDBforMysql) ConnectDB() (dbconn *sql.DB, err error) {

	port_list := strconv.Itoa(ty.stDBConnectInfo.NPort)

	stringArray := []string{ty.stDBConnectInfo.StrID, ":", ty.stDBConnectInfo.StrPasswd, "@tcp(", ty.stDBConnectInfo.StrIP, ":", port_list, ")/", ty.stDBConnectInfo.StrDBname}
	justString := strings.Join(stringArray, "")

	db, err := sql.Open("mysql", justString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(ty.mMaxIdelConn)
	db.SetMaxOpenConns(ty.mMaxOpenConn)

	ty.mDBconn = db
	ty.bInited = true

	ty.mJobTime = time.Now()
	ty.closeFlg = false
	go ty.dbIdleUnconnectionChecker()
	return db, nil
}

func (ty *BCEDBforMysql) DBClose() {
	if ty.mDBconn != nil {
		ty.mDBconn.Close()
		ty.closeFlg = true
	}
}

//dbIdleUnconnectionChecker 특정시간 이상 디비작업 없을경우 mysql이 커넥션을 끊는다. 이를 해결하기위해 마지막 작업타임을 저장후 특정시간 이상 작업이 없다면 의미없는 작업을 진행하여 연결유지하도록 한다.
func (ty *BCEDBforMysql) dbIdleUnconnectionChecker() {
	for {
		if ty.closeFlg == true {
			return
		}

		time.Sleep(time.Duration(int64(time.Second)))

		now := time.Now()

		idT := ty.mJobTime.Add(time.Duration(int64(time.Second) * ty.idleTime))

		if idT.Unix() < now.Unix() {
			sql := "select"
			ty.DBQuerySelect(sql)
		}
	}

}

//GameDBQuerySelect 게임디비 셀렉트
func (ty *BCEDBforMysql) DBQuerySelect(query string) (r map[int]map[string]interface{}, err error) {
	db := ty.mDBconn
	ty.mJobTime = time.Now()
	if db == nil {
		return nil, fmt.Errorf("Database not connection")
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	send := make(map[int]map[string]interface{}, 13)
	i := 0

	for rows.Next() {
		tmp := make(map[string]interface{})
		row := make([]interface{}, len(columns))
		for idx := range columns {
			row[idx] = new(MetalScanner)
		}

		err := rows.Scan(row...)
		if err != nil {
			if rows != nil {
				rows.Close()
			}
			return nil, err
		}

		for idx, column := range columns {
			var scanner = row[idx].(*MetalScanner)
			strtmp := scanner.value
			tmp[column] = strtmp
		}
		send[i] = tmp
		i++
	}
	if rows != nil {
		rows.Close()
	}
	// 순서랑은 상관없이 진행될 수 이씅ㅁ
	return send, nil
}

//GameDBQueryExec 쿼리 실행(Insert Update Delete등)
func (ty *BCEDBforMysql) DBQueryExec(query string) (nAffe int64, er error) {
	if ty.stDBConnectInfo.Readonly {
		return 0, fmt.Errorf("db select only")
	}
	db := ty.mDBconn
	ty.mJobTime = time.Now()
	if db == nil {
		return 0, fmt.Errorf("Database not connection")
	}

	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (ty *BCEDBforMysql) DBQueryExecInterface(query string, args ...interface{}) (nAffe int64, er error) {
	if ty.stDBConnectInfo.Readonly {
		return 0, fmt.Errorf("db select only")
	}

	db := ty.mDBconn
	ty.mJobTime = time.Now()
	if db == nil {
		return 0, fmt.Errorf("Database not connection")
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	// sql.Result.RowsAffected() 체크
	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (ty *BCEDBforMysql) GetErrorCode(str string) int {

	//Error 1062: Duplicate entry '34' for key 'PRIMARY'
	index := strings.IndexAny(str, ":")
	result := str[5:index]
	response := strings.Join(strings.Fields(result), "")
	code, err := strconv.Atoi(response)
	if err != nil {
		return -1
	}
	return code
}

func (scanner *MetalScanner) getBytes(src interface{}) []byte {
	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}

func (scanner *MetalScanner) Scan(src interface{}) error {
	switch src.(type) {
	case int:
		if value, ok := src.(int); ok {
			scanner.value = value
			scanner.valid = true
		}
	case int64:
		if value, ok := src.(int64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case float64:
		if value, ok := src.(float64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case float32:
		if value, ok := src.(float32); ok {
			scanner.value = value
			scanner.valid = true
		}
	case bool:
		if value, ok := src.(bool); ok {
			scanner.value = value
			scanner.valid = true
		}
	case string:
		value := scanner.getBytes(src)
		scanner.value = string(value)
		scanner.valid = true
	case []byte:
		value := scanner.getBytes(src)
		scanner.value = string(value)
		scanner.valid = true
	case time.Time:
		if value, ok := src.(time.Time); ok {
			scanner.value = value
			scanner.valid = true
		}
	case nil:
		scanner.value = nil
		scanner.valid = true
	}
	return nil
}

// sql 만들기 힘드신 분들이 사용하는 struct 를 이용한 DB insert

func formatString(arg interface{}) string {
	switch arg.(type) {
	case int:
		i := arg.(int)
		return strconv.Itoa(i)
	case int64:
		i := arg.(int64)
		return strconv.FormatInt(i, 10)
	case float32:
		f := arg.(float32)
		return strconv.FormatFloat(float64(f), 'f', -1, 32)
	case float64:
		f := arg.(float64)
		return strconv.FormatFloat(f, 'f', -1, 64)
	case string:
		s := arg.(string)
		return s
	default:
		return "Error"
	}
}

func FormatSqlInsert(tableName string, coulumeText []string, arg ...interface{}) string {
	// string
	/*
		sql := fmt.Sprintf("insert into %s(
		)" )
	*/

	sql := fmt.Sprintf("insert into %s(", tableName)

	for _, value := range coulumeText {
		sql += "`" + value + "`,"
	}
	sql = string(sql[0:(len(sql) - 1)])
	sql += ") values ("

	for _, v := range arg {
		switch v.(type) {
		case string:
			sql += "'" + formatString(v) + "',"
		default:
			sql = sql + formatString(v) + ","
		}
	}
	sql = string(sql[0:(len(sql) - 1)])
	sql = sql + ")"

	return sql

}

func SturctToMap(tagName string, obj interface{}) map[string]interface{} {
	temp := structs.New(obj)
	mapToFace := structs.Map(obj)
	makeMap := make(map[string]interface{})

	for k, v := range mapToFace {
		//name := s.Field(k)
		tags := temp.Field(k)
		key := tags.Tag(tagName)
		if key != "" {
			makeMap[key] = v
		}
	}
	return makeMap
}

func MapToSqlInsert(tablesName string, maps map[string]interface{}) string {

	var columStr []string
	var values []interface{}

	for k, v := range maps {
		columStr = append(columStr, k)
		values = append(values, v)
	}

	sql := FormatSqlInsert(tablesName, columStr, values...)
	return sql

}
