package FP_DAO

import (
	"FP_DAO/FP_DB"
	"FP_Model"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	_SsThreshold = 1
	_ChannelBuff = 10
	_WorkerNr    = 5
)

type userDAO struct{}

func NewUserDAO() IUserDAO {
	return &userDAO{}
}

func (this *userDAO) GetUserByUid(uid int64) (u *FP_Model.User, e error) {
	rows, res := FP_DB.GetUserByUId(uid)
	if len(rows) == 0 {
		u, e = nil, NotSuchUidError
	} else {
		u, e = parseUser(uid, rows, res), nil
	}
	return
}

//Send all uid by the channel
func uidChannelSender(limit int) <-chan int64 {
	uidChan := make(chan int64, _ChannelBuff)

	go func(ch chan<- int64) {
		for _, uid := range FP_DB.GetAllUserUid(limit) {
			ch <- uid
		}
		close(ch) //Close by sender
	}(uidChan)

	return uidChan
}

func (this *userDAO) GetAllUsers(limit int) []*FP_Model.User {
	if limit < 0 {
		limit = 1
	}
	users := make([]*FP_Model.User, 0, limit)
	mutex := sync.Mutex{}
	uidChan := uidChannelSender(limit)
	signal := make(chan bool)
	//Send workers
	for ii := 0; ii < _WorkerNr; ii++ {
		go func(ch <-chan int64, signal chan<- bool) {

			for uid := range ch {
				if u, err := this.GetUserByUid(uid); err == nil {
					if u.SS() > _SsThreshold {
						mutex.Lock()
						users = append(users, u)
						mutex.Unlock()
					} else {
						log.Printf("User :%d Not enough information.\n", u.Uid)
					}
				}
			}

			signal <- true
		}(uidChan, signal)
	}

	for ii := 0; ii < _WorkerNr; ii++ {
		<-signal
	}
	return users
}

//Parse the result from the database into the USER model
func parseUser(uid int64, rows []mysql.Row, res mysql.Result) *FP_Model.User {
	user := FP_Model.NewUser()
	user.Uid = uid

	txtIdx := res.Map("txt")
	levelIdx := res.Map("level")

	for _, row := range rows {
		txt := row.Str(txtIdx)
		level := row.Str(levelIdx)
		parseTxt(txt, level, user)
	}

	return user
}

func parseTxt(txt string, level string, user *FP_Model.User) {
	factor := 1
	switch level {
	case "NB":
		factor = 16
	case "JL":
		factor = 4
	}

	tfs := strings.Split(txt, "$")
	for _, tf := range tfs {
		fields := strings.Fields(tf)
		if len(fields) != 2 {
			continue
		}
		if weight, err := strconv.Atoi(fields[1]); err == nil {
			var value int32 = int32(weight * factor)
			//Incase for overflow
			if value > 5200 {
				value = 5200
			}
			user.Put(fields[0], value)
		}
	}
}
