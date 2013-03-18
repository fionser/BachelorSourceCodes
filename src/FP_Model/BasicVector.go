package FP_Model

import (
	"math"
)

//BasicVector implements the interface IVector.
//BasicVector is embeded by User and Center.
type BasicVector struct {
	///Term-Frenquency
	_TF map[string]int32
	//Sum of square of each dimensions.
	//To accelerate the cosine calculation between vectors.
	_SS int32
}

func newBasicVector() *BasicVector {
	v := new(BasicVector)
	v.constructor()
	return v
}

//Return the cosine angle between vector u and v.
func cosine(u, v IVector) float64 {
	//cos<u, v> = u*v / (|u| * |v|)
	var fraction int32 = 0
	var factor float64 = math.Sqrt(float64(u.SS())) * math.Sqrt(float64(v.SS()))

	for k := range u.Iterator() {
		if val2, ok := v.TryGet(k); ok {
			fraction += (u.Get(k) * val2)
		}
	}

	if factor-float64(fraction) < 1e-8 {
		return 1
	}
	return float64(fraction) / factor
}

//Use cosine angle as the measure of distance 
func cosineAsDitance(u, v IVector) float64 {
	return math.Acos(cosine(u, v))
}

func (this *BasicVector) constructor() {
	this._TF = make(map[string]int32)
	this._SS = 0
}

func euclideanAsDistance(u, v IVector) float64 {
	d := float64(u.SS()) + float64(v.SS())
	return math.Sqrt(d)
}

//Return the distance between vector this and u.
func (this *BasicVector) Distance(u IVector) float64 {
	return cosineAsDitance(this, u)
}

func (this *BasicVector) EDistance(u IVector) float64 {
	return euclideanAsDistance(this, u)
}

//Insert a key-value pair into the vector.
//Will modify the _SS value.
func (this *BasicVector) Put(key string, val int32) {
	//Substract the old value from _SS if a renew key is inserted.
	if oldVal, ok := this._TF[key]; ok {
		this._SS -= oldVal * oldVal
	}
	//Add the square of the new value into the _SS.
	this._TF[key] = val
	//In case for overflow
	if math.MaxInt32-this._SS > val*val {
		this._SS += val * val
	} else {
		this._SS = math.MaxInt32
	}
}

//Reset the vector to zero.
func (this *BasicVector) CleanUp() {
	//this._TF = nil
	//this._TF = make(map[string]int32)
	for k, _ := range this._TF {
		delete(this._TF, k)
	}
	this._SS = 0
}

//Try to get the value with the key
//return false if the key is unfounded.
func (this *BasicVector) TryGet(key string) (val int32, ok bool) {
	val, ok = this._TF[key]
	return
}

//To get the value with the key for sure.
//Panic with throw if key is unfounded.
func (this *BasicVector) Get(key string) (val int32) {
	if v, ok := this._TF[key]; !ok {
		panic("Mismatch key")
	} else {
		val = v
	}
	return
}

//Return the iterator of all the key in User, for the useage in range.
//Channel will be closed after all key is sent.
func (this *BasicVector) Iterator() <-chan string {
	ch := make(chan string, 10)
	go func() {
		for k, _ := range this._TF {
			ch <- k
		}
		close(ch)
	}()
	return ch
}

//Return the sum of square 
func (this *BasicVector) SS() (SS int32) {
	//In case for the zero factor in function cosine
	if this._SS == 0 {
		SS = 1
	} else if this._SS < 0 {
		panic("Negative SS")
	} else {
		SS = this._SS
	}
	return
}
