package FP_Logic

import (
	"FP_Domain"
	"FP_Util"
)

var (
	TRY_TIME int
)

func init() {
	TRY_TIME = FP_Util.GetInstance().GetInt("TRY_TIME", 10)
}

type BinaryAggregator struct {
	*_KMeanAggregator
}

func NewBinaryAggregator() *BinaryAggregator {
	ba := &BinaryAggregator{}
	ba.constructor()
	return ba
}

func (this *BinaryAggregator) constructor() {
	this._KMeanAggregator = NewKMeanAggregator(1)
}

func (this *BinaryAggregator) IsFinish() bool {
	if len(this.clusters) == _CLUSTER_NR {
		return true
	}
	return false
}

func (this *BinaryAggregator) getMaxESCluster() *FP_Domain.Cluster {
	var maxEs float64 = -1
	var maxC *FP_Domain.Cluster
	var maxIdx int
	for idx, c := range this.clusters {
		es := c.ES()
		if es > maxEs {
			maxIdx = idx
			maxC = c
			maxEs = es
		}
	}
	//Remove the max one.
	this.clusters = append(this.clusters[0:maxIdx], this.clusters[maxIdx+1:]...)
	return maxC
}

func depyCopy(clusters []*FP_Domain.Cluster) []*FP_Domain.Cluster {
	clts := make([]*FP_Domain.Cluster, len(clusters), len(clusters))

	for ii := 0; ii < len(clusters); ii++ {
		clts[ii] = FP_Domain.CopyCluster(clusters[ii])
	}
	return clts
}

func (this *BinaryAggregator) addSplitted(clts []*FP_Domain.Cluster) {
	for _, c := range clts {
		if c.Size() > 0 {
			this.clusters = append(this.clusters, c)
		}
	}
}

func (this *BinaryAggregator) InitClusters() {
	this.clusters[0] = FP_Domain.NewCluster(this.userSet.Index(0), -1)
	for u := range this.userSet.Iterator() {
		this.clusters[0].Add(u)
	}
}

func (this *BinaryAggregator) Aggregate() {
	this.InitUserSet()
	this.InitClusters()
	kMeanAggregator := NewKMeanAggregator(2)

	for {
		c := this.getMaxESCluster()
		println("SIZE ", c.Size())
		kMeanAggregator.userSet = c.UserSet
		kMeanAggregator.InitClusters()
		kMeanAggregator._SES = 0
		kMeanAggregator.first = true

		var try_time = 0
		for !kMeanAggregator.IsFinish() && try_time < TRY_TIME {
			kMeanAggregator.cleanClusters()
			kMeanAggregator.AssignUser()
			kMeanAggregator.GetNewCenter()
			try_time++
		}
		size1 := kMeanAggregator.clusters[0].Size()
		size2 := kMeanAggregator.clusters[1].Size()
		println("SPLIT ", size1, size2)
		this.addSplitted(depyCopy(kMeanAggregator.clusters))

		if this.IsFinish() {
			break
		}
	}

	this.printResult()
}

func (this *BinaryAggregator) AverageRadius() float64 {
	var radius float64 = 0
	var size = len(this.clusters)
	for ii := 0; ii < size; ii++ {
		radius += this.clusters[ii].Radius()
	}
	return radius / float64(size)
}
