package FP_Logic

import (
	"FP_Domain"
	"FP_Model"
	"FP_Util"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	//Threshold for aggregation iteration end-guard.
	_THRESHOLD = 0.0001
	_WORKER_NR = 5
)

var (
	_CLUSTER_NR   int
	_USERSET_SIZE int
)

func init() {
	configer := FP_Util.GetInstance()
	_CLUSTER_NR = configer.GetInt("CLUSTER_NR", 46)
	_USERSET_SIZE = configer.GetInt("USERSET_SIZE", 2000)
}

type _KMeanAggregator struct {
	userSet        *FP_Domain.UserSet
	clusters       []*FP_Domain.Cluster
	first          bool
	_SES           float64 // Sum of Error Square
	distanceMatrix [][]float64
}

func NewKMeanAggregator(clusterNr int) *_KMeanAggregator {
	ka := &_KMeanAggregator{}
	if clusterNr <= 0 {
		clusterNr = _CLUSTER_NR
	}
	ka.constructor(clusterNr)
	return ka
}

func (this *_KMeanAggregator) constructor(initClusterNr int) {
	this.clusters = make([]*FP_Domain.Cluster, initClusterNr, _CLUSTER_NR)
	this.distanceMatrix = make([][]float64, _CLUSTER_NR)
	for ii := range this.distanceMatrix {
		//Down-Left triangle Only
		this.distanceMatrix[ii] = make([]float64, ii+1)
	}
	this._SES = 0
	this.first = true
}

func (this *_KMeanAggregator) printResult() {
	configer := FP_Util.GetInstance()
	output := configer.GetString("OUTPUT", "/tmp/out.txt")
	fd, err := os.Create(output)
	var writer io.Writer

	if err == nil {
		defer fd.Close()
		writer = fd
	} else {
		log.Println(err.Error())
		writer = os.Stdout
	}

	for idx, c := range this.clusters {
		fmt.Fprintf(writer, "Cluster %d. Size %d\n", idx, c.Size())
		for u := range c.Iterator() {
			fmt.Fprintf(writer, "%d\n", u.Uid)
		}
	}
}

func (this *_KMeanAggregator) Aggregate() {
	log.Println("Going to aggregate.")
	this.InitUserSet()
	this.InitClusters()
	for !this.IsFinish() {
		this.cleanClusters()
		this.AssignUser()
		this.GetNewCenter()
	}
	this.printResult()
}

//Randomly select the original centers for the beginning clusters
func (this *_KMeanAggregator) _RandomlySelectCenters() {
	nr := len(this.clusters)
	randomSelector :=
		FP_Util.NewSelector(int32(nr), this.userSet.Size())

	for i := 0; i < nr; i++ {
		idx := randomSelector.Next()
		this.clusters[i] = FP_Domain.NewCluster(this.userSet.Index(idx), -1)
	}
}

//Select the centers which have the maximum sum of distance between each other.
func (this *_KMeanAggregator) _MaxDistanceCenters() {
	if len(this.clusters) == 2 {
		rand.Seed(time.Now().Unix())
		firstCenter := this.userSet.Index(rand.Int31() % this.userSet.Size())
		var maxDist float64 = -1
		var secondCenter *FP_Model.User
		for u := range this.userSet.Iterator() {
			dist := u.Distance(firstCenter)
			if dist >= maxDist {
				secondCenter = u
				maxDist = dist
			}
		}

		this.clusters[0] = FP_Domain.NewCluster(firstCenter, -1)
		this.clusters[1] = FP_Domain.NewCluster(secondCenter, -1)
	} else {
		this._RandomlySelectCenters()
	}
}

func (this *_KMeanAggregator) calDistanceBetweenClusters() {
	size := len(this.clusters)
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			this.distanceMatrix[i][j] =
				this.clusters[i].Center().Distance(this.clusters[j].Center())
		}
	}
}

func (this *_KMeanAggregator) InitClusters() {
	//log.Println("Going to initialize clusters")
	//this._MaxDistanceCenters()
	this._RandomlySelectCenters()
	//this.calDistanceBetweenClusters()
	//log.Println("Clusters initialized.")
}

func (this *_KMeanAggregator) InitUserSet() {
	//log.Println("Going to initialize userset")
	this.userSet = FP_Domain.LoadFromDB(_USERSET_SIZE)
	log.Println("UserSet initialized. Size: ", this.userSet.Size())
}

func (this *_KMeanAggregator) getDistanceFromMatrix(i, j int) float64 {
	var dist float64
	if i > j {
		dist = this.distanceMatrix[i][j]
	} else {
		dist = this.distanceMatrix[j][i]
	}
	return dist
}

func (this *_KMeanAggregator) getBiggestClusterIndex() int {
	var maxSize int = -1
	var maxIndex int
	for ii, c := range this.clusters {
		size := int(c.Size())
		if size > maxSize {
			maxIndex = ii
			maxSize = size
		}
	}
	return maxIndex
}

func (this *_KMeanAggregator) assignUser(signal chan<- bool,
	userIterator <-chan *FP_Model.User) {
	var closestCluster *FP_Domain.Cluster
	var clustersSize int = len(this.clusters)
	//For each user in the userset
	//to find the closest cluster and insert the user into
	//that cluster
	for u := range userIterator {
		var dist float64 = math.MaxFloat64
		for i := 0; i < clustersSize; i++ {
			tmp := u.Distance(this.clusters[i].Center())
			if tmp < dist {
				closestCluster = this.clusters[i]
				dist = tmp
			}
		}
		////Assume that User 'u' is closest to biggest cluster
		//firstIndex := this.getBiggestClusterIndex()
		//minIndex := firstIndex
		//closestCluster = this.clusters[minIndex]
		//minDist = u.Distance(this.clusters[minIndex].Center())
		//for i := 0; i < clustersSize; i++ {
		//	if i == firstIndex {
		//		continue
		//	}
		//	//According Triangle Inequality: d(x, z) <= d(x, y) + d(y, z), we have:
		//	// 2d(x, c_i) <= d(c_i, c_k) ==> d(x, c_i) <= d(x, c_k)
		//	//PS. 循环不变式: i >= minIndex
		//	if 2*minDist <= this.getDistanceFromMatrix(i, minIndex) {
		//		continue
		//	} else {
		//		dist := u.Distance(this.clusters[i].Center())
		//		if dist <= minDist {
		//			//equals, 50% assigning to new cluster.
		//			if dist == minDist {
		//				if rand.Int()&1 == 0 {
		//					continue
		//				}
		//			}
		//			minIndex = i
		//			minDist = dist
		//			closestCluster = this.clusters[i]
		//		}
		//	}
		//}
		//Assign the user to the closest cluster
		closestCluster.Add(u)
		//closestCluster.NewCenter()
	}
	//Signal the caller its finishing.
	signal <- true
}

func (this *_KMeanAggregator) AssignUser() {
	//log.Println("Going to assign user.")
	signal := make(chan bool, _WORKER_NR)
	userIterator := this.userSet.Iterator()

	for i := 0; i < _WORKER_NR; i++ {
		//Send worker goroutines
		go this.assignUser(signal, userIterator)
	}
	//Waiting workers to finish.
	for i := 0; i < _WORKER_NR; i++ {
		<-signal
	}
	log.Println("User assignment finished.")
}

func (this *_KMeanAggregator) GetNewCenter() {
	for _, c := range this.clusters {
		c.NewCenter()
	}
	//Renew the distance metrix
	this.calDistanceBetweenClusters()
	log.Println("New center finished.")
}

func (this *_KMeanAggregator) cleanClusters() {
	for _, c := range this.clusters {
		c.RemoveAll()
	}
}

//Get the sum of error square from these clusters.
func (this *_KMeanAggregator) GetSES() float64 {
	log.Println("Going to get SES.")
	var ses float64 = 0
	cnt := len(this.clusters)
	ch := make(chan float64, cnt)

	for i, _ := range this.clusters {
		go func(ch chan<- float64, c *FP_Domain.Cluster) {
			ch <- c.ES()
		}(ch, this.clusters[i])
		/*  //Code below don't works.
		go func(ch chan<- float64) {
			ch <- c.ES()
		}(ch)
		*/
	}

	for ii := 0; ii < cnt; ii++ {
		ses += <-ch
	}

	close(ch)
	return ses
}

func (this *_KMeanAggregator) IsFinish() bool {
	//First iteration.
	if this.first {
		this.first = false
		return false
	}

	ses := this.GetSES()
	log.Printf("Old ses %.4f. New ses %.4f\n", this._SES, ses)
	if ses == 0 {
		for _, c := range this.clusters {
			println(c.Size())
		}
	}

	// The change is less than the threshold.
	if math.Abs(ses-this._SES) < this._SES*_THRESHOLD {
		log.Println("KMean Finish")
		return true
	}

	this._SES = ses
	return false
}

func (this *_KMeanAggregator) AverageRadius() float64 {
	var radius float64 = 0
	for _, c := range this.clusters {
		radius += c.Radius()
	}
	return radius / float64(len(this.clusters))
}
