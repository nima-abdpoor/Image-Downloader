package model

type ImageResult struct {
	Id      int
	QueryId string
	Url     string
	Data    string
}

type SearchImage struct {
	Query  string
	Status string
	Date   string
}

const (
	StatusInProgress = "InProgress"
	StatusFailed     = "Failed"
	StatusSuccess    = "Successful"
)
