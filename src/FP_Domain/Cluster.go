package FP_Domain

import (
	"FP_Model"
	"math"
	"sync"
)

const (
	DEFAULT_CLUSTER_SIZE = 2000
	_ZERO                = 10e-5
)

type Cluster struct {
	center *FP_Model.Center
	mutex  sync.Mutex
	*UserSet
}

//Create a new cluster.
//Need a vector as the original center.
//Default size will be used when ClusterSize is less than 1.
func NewCluster(v FP_Model.IVector, ClusterSize int) *Cluster {
	c := new(Cluster)
	c.constructor(v, ClusterSize)
	return c
}

func CopyCluster(c *Cluster) *Cluster {
	nc := NewCluster(c.center, int(c.Size()))
	for u := range c.Iterator() {
		nc.Add(u)
	}
	return nc
}

func (this *Cluster) constructor(c FP_Model.IVector, ClusterSize int) {
	this.center = FP_Model.NewCenter(c)

	if ClusterSize < 1 {
		ClusterSize = DEFAULT_CLUSTER_SIZE
	}
	this.UserSet = NewUserSet(ClusterSize)
}

//Calculate the new center of the cluster
func (this *Cluster) NewCenter() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.center.CleanUp()
	//Iterate users in cluster
	for u := range this.Iterator() {
		//Iterate keys in user
		for key := range u.Iterator() {
			v := u.Get(key)
			if v1, ok := this.center.TryGet(key); ok {
				this.center.Put(key, v1+v)
			} else {
				this.center.Put(key, v)
			}
		}
	}

	this.center.Means(this.Size())
}

//Get the error square of the cluster
//Using the sum of distance to center as the ES.
func (this *Cluster) ES() float64 {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var es float64 = 0
	for u := range this.Iterator() {
		e := math.Pow(u.Distance(this.center), 2.0)
		es += e
	}

	return es
}

//The radius of the cluster is the maximum distance between
//vectors and the center.
func (this *Cluster) Radius() float64 {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var radius float64 = 0
	var tmp float64
	for u := range this.Iterator() {
		tmp = this.center.EDistance(u)
		if tmp > radius {
			radius = tmp
		}
	}
	return radius
}

func (this *Cluster) Center() *FP_Model.Center {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.center
}
