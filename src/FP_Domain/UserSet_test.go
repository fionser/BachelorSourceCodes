package FP_Domain

import (
	"FP_Model"
	"testing"
)

func TestCluster(t *testing.T) {
	c := FP_Model.NewCenter(nil)
	c.Put("1", 1.0)
	c.Put("2", 2.0)
	c.Put("3", 3.0)
	cluster := NewCluster(c, 10)
	for ii := 1; ii <= 3; ii++ {
		nc := FP_Model.NewUser()
		va := int32(ii)
		nc.Put("1", va)
		nc.Put("2", va)
		nc.Put("3", va)
		cluster.Add(nc)
	}
	println(cluster.ES())
	cluster.NewCenter()
	println(cluster.ES())
}
