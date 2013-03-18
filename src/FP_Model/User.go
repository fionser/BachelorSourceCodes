package FP_Model

type User struct {
	Uid int64
	//Embed the interface rather than the
	//implementation.
	IVector
}

func NewUser() *User {
	u := new(User)
	u.constructor()
	return u
}

func (this *User) constructor() {
	this.Uid = -1
	//We can switch different implementation here.
	//but not change the interface.
	this.IVector = newBasicVector()
}
