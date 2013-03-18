package FP_Util

import (
	"math/rand"
	"time"
)

type RandomSelector struct {
	_nr, _range int32
	_prev       int32
}

//Randomly select nr integers from range [1...rge]
func NewSelector(nr, rge int32) *RandomSelector {
	selector := &RandomSelector{nr, rge, 0}
	return selector
}

func (this *RandomSelector) Reset(nr, rge int32) {
	this._nr = nr
	this._range = rge
	this._prev = 0
}

func (this *RandomSelector) Next() (nxt int32) {
	rand.Seed(time.Now().Unix())

	for ; this._prev < this._range; this._prev++ {
		if rand.Int31()%(this._range-this._prev) < this._nr {
			nxt = this._prev
			this._prev++
			this._nr--
			break
		}
	}
	return
}
