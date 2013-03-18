package FP_DB

import (
	"FP_Util"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"os"
	"sync"
)

var (
	user     string
	pass     string
	dbnm     string
	prot     string
	addr     string
	db       mysql.Conn
	stmtLock sync.Mutex
	uidStmt  mysql.Stmt
)

func init() {
	loadConfig()
	db = mysql.New(prot, "", addr, user, pass, dbnm)
	checkError(db.Connect())
	uidStmt, _ = db.Prepare("select txt, level from TF where uid = ?")
}

func loadConfig() {
	configer := FP_Util.GetInstance()
	user = configer.GetString("DB_USER", "root")
	pass = configer.GetString("DB_PASS", "44586218")
	dbnm = configer.GetString("DB_DBNM", "friends_recommendation")
	prot = configer.GetString("DB_PROT", "tcp")
	addr = configer.GetString("DB_ADDR", "localhost:3306")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkResult(rows []mysql.Row, res mysql.Result, err error) ([]mysql.Row, mysql.Result) {
	checkError(err)
	return rows, res
}

//Return row SQL result from database
func GetUserByUId(uid int64) ([]mysql.Row, mysql.Result) {
	stmtLock.Lock()
	defer stmtLock.Unlock()
	uidStmt.Bind(uid)
	return checkResult(uidStmt.Exec())
}

//Get All Uid from database.Limition is needed.
func GetAllUserUid(limit int) []int64 {
	if limit < 0 {
		limit = 0
	}
	uids := make([]int64, limit, limit)
	stmt, _ := db.Prepare("select distinct(uid) from TF limit ?")
	rows, res := checkResult(stmt.Exec(limit))
	uidIdx := res.Map("uid")
	for idx, row := range rows {
		uids[idx] = row.Int64(uidIdx)
	}
	return uids
}
