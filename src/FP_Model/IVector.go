package FP_Model

//Definition of basic vectorical operation.
type IVector interface {
	//Return the distance between vector this and u.
	Distance(u IVector) float64
	//Return the euclidean distance
	EDistance(u IVector) float64
	//Insert a key-value pair into the vector.
	//Will modify the ss value.
	Put(key string, value int32)
	//To get the value with the key for sure.
	//Panic with throw if key is unfounded
	//Cooperate with function Iterator.
	Get(key string) int32
	//Try to get the value with the key
	//return false if the key is unfounded.
	TryGet(key string) (int32, bool)
	//Return the iterator of all the key in vector.
	//Channel will be closed after all key is sent.
	Iterator() <-chan string
	//Return the sum of square 
	SS() int32
	//Reset the vector to zero.
	CleanUp()
}
