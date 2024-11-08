package repository

import (
	"GoogleImageDownloader/model"
	"sync"
)

type SearchEngine interface {
	SearchImages(query string, page, perPage int32, wg *sync.WaitGroup) ([]model.ImageResult, error)
}
