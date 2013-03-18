package FP_DAO

import (
	"FP_Model"
	"testing"
)

func printUser(u *FP_Model.User) {
	if u == nil {
		return
	}
	for k, v := range u.TF {
		print(k, v)
	}
	println("\n********")
}

func TestUserDAO(t *testing.T) {
	userdao := NewUserDAO()
	users := userdao.GetAllUsers(3)
	println(len(users))
}
