package FP_Logic

// K-Mean Aggregator 
// Algorithm steps as fellow:
// Initialize the clusters and userset.
// Assign all users into each cluster.
// Calculate the new center of each cluster.
// Go into next iteration if is not finish.
type IAggregator interface {
	Aggregate()
	InitClusters()
	InitUserSet()
	AssignUser()
	GetNewCenter()
	IsFinish() bool
	GetSES() float64
	AverageRadius() float64
}

func NewAggregator() IAggregator {
	//agg := NewKMeanAggregator(-1)
	return NewBinaryAggregator()
}
