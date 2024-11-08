package service

import (
	"GoogleImageDownloader/repository"
	"sync"
)

type SearchImage interface {
	SearchImage(searchEngine repository.SearchEngine, query string, number, page int32, wg *sync.WaitGroup)
	RetrySearchImage(searchEngine repository.SearchEngine, query string, queryId int64, number, page int32, wg *sync.WaitGroup)
}
