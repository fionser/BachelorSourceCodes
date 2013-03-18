package FP_DAO

import (
	"FP_Model"
	"errors"
)

var (
	NotSuchUidError = errors.New("NotSuchUidError")
)

//UserDAO interface
type IUserDAO interface {
	//NotSuchUidError will be returned if uid is unfounded.
	GetUserByUid(uid int64) (*FP_Model.User, error)
	//Get Users from database.Limition is needed.
	GetAllUsers(limit int) []*FP_Model.User
}
