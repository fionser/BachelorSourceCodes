package FP_Domain

import (
	"FP_DAO"
	"FP_Model"
	"sync"
)

const (
	DEFAULT_SET_SIZE = 3000
	_ITERATOR_BUFFER = 10
)

type UserSet struct {
	users []*FP_Model.User
	mutex sync.Mutex
}

func (this *UserSet) constructor(SetSize int) {
	if SetSize < 1 {
		SetSize = DEFAULT_SET_SIZE
	}
	this.users = make([]*FP_Model.User, 0, SetSize)
}

//Create a UserSet
//Default size will be used if the parameter SetSize is less than 1.
func NewUserSet(SetSize int) *UserSet {
	us := new(UserSet)
	us.constructor(SetSize)
	return us
}

//Load a UserSet from Database.
//Parameter limit is the maximum number users to load.  
func LoadFromDB(limit int) *UserSet {
	if limit < 0 {
		limit = 1
	}
	userSet := new(UserSet)
	//Just asign one slot for the user.
	//Because we will replace the content later.
	userSet.constructor(1)

	userDAO := FP_DAO.NewUserDAO()
	userSet.users = userDAO.GetAllUsers(limit)

	return userSet
}

//Thread-safely adding a new user into
//the UserSet.
//Nil parameter will be ignored.
func (this *UserSet) Add(u *FP_Model.User) {
	if u != nil {
		this.mutex.Lock()
		defer this.mutex.Unlock()
		this.users = append(this.users, u)
	}
}

//Iterate the users in the UserSet by index
func (this *UserSet) Index(idx int32) *FP_Model.User {
	//this.mutex.Lock()
	//defer this.mutex.Unlock()
	if idx >= 0 && idx < this.Size() {
		return this.users[idx]
	}
	return nil
}

//Thread-safe
//Return the number of users in the UserSet
func (this *UserSet) Size() int32 {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return int32(len(this.users))
}

func (this *UserSet) RemoveAll() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.users = this.users[:0]
}

//Lock the whole userset when iterating. Unlock after the iteration finished.
func (this *UserSet) Iterator() <-chan *FP_Model.User {
	this.mutex.Lock()
	ch := make(chan *FP_Model.User, _ITERATOR_BUFFER)

	go func() {
		for _, u := range this.users {
			ch <- u
		}
		close(ch)
		this.mutex.Unlock()
	}()

	return ch
}
