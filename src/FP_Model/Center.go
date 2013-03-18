package FP_Model

//Center for Cluster
type Center struct {
	IVector
}

//If parameter v is nil, then just return an
//empty center,otherwise deeply copy the
//parameter v into this center
func NewCenter(v IVector) *Center {
	c := new(Center)
	c.constructor(v)
	return c
}

func (this *Center) constructor(v IVector) {
	this.IVector = newBasicVector()
	if v != nil {
		for k := range v.Iterator() {
			this.Put(k, v.Get(k))
		}
	}
}

//To mean each dimension by parameter k in this center.
//Dimension will be low-rounded
func (this *Center) Means(k int32) {
	if k < 0 {
		panic("illegal parameter")
	}
	if k > 0 {
		for key := range this.Iterator() {
			this.Put(key, this.Get(key)/k)
		}
	}
}
