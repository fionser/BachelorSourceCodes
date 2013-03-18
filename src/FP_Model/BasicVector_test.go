package FP_Model

import (
	"testing"
)

func TestBV(t *testing.T) {
	v1 := NewUser()
	v1.Put("1", 1.0)
	c := NewCenter(v1)
	v2 := NewUser()
	v2.Put("1", 1.0)
	println(c.Distance(v2))
}
