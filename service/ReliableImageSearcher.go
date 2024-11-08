package service

import (
	db "GoogleImageDownloader/db/sqlc"
	"GoogleImageDownloader/model"
	"GoogleImageDownloader/repository"
	"context"
	"fmt"
	"sync"
	"time"
)

type ReliableImageSearcher struct {
	Store  *db.Store
	config Config
}

type Config struct {
	TimeOut time.Duration `koanf:"time_out"`
}

func New(store *db.Store, config Config) SearchImage {
	return ReliableImageSearcher{
		Store:  store,
		config: config,
	}
}

func (searcher ReliableImageSearcher) SearchImage(searchEngine repository.SearchEngine, query string, number, page int32, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), searcher.config.TimeOut)
	defer cancel()
	queryResult, err := searcher.Store.CreateQuery(ctx, query, number, page)
	if err != nil {
		fmt.Println("error in Inserting Into Query!", err)
	}

	if imageResults, err := searchEngine.SearchImages(query, page, number, wg); err != nil {
		fmt.Printf("Error in Searching Images %v", err)

		dErr := searcher.Store.UpdateQuery(context.Background(), db.UpdateQueryArgs{
			Id:     queryResult.Id,
			Status: model.StatusFailed,
		})
		if dErr != nil {
			fmt.Printf("Error in Updaing Failed Request %v", dErr)
		}
	} else {
		for _, result := range imageResults {
			err := searcher.Store.CreateImageResult(ctx, db.CreateImageParams{
				QueryID: queryResult.Id,
				Url:     result.Url,
				Data:    result.Data,
			})
			if err != nil {
				fmt.Printf("Error In storing ImageResult with url:%v err: %v", result.Url, err)
			}
		}

		fmt.Printf("RESULT: query:%s, len:%d page:%d perPage:%d\n", query, len(imageResults), page, number)

		dErr := searcher.Store.UpdateQuery(context.Background(), db.UpdateQueryArgs{
			Id:     queryResult.Id,
			Status: model.StatusSuccess,
		})
		if dErr != nil {
			fmt.Printf("Error in Updaing Success Request %v", dErr)
		}
	}
}

func (searcher ReliableImageSearcher) RetrySearchImage(searchEngine repository.SearchEngine, query string, queryId int64, number, page int32, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), searcher.config.TimeOut)
	defer cancel()

	if imageResults, err := searchEngine.SearchImages(query, page, number, wg); err != nil {
		fmt.Printf("Error in Searching Images %v", err)

		dErr := searcher.Store.UpdateQuery(context.Background(), db.UpdateQueryArgs{
			Id:     queryId,
			Status: model.StatusFailed,
		})
		if dErr != nil {
			fmt.Printf("Error in Updaing Failed Request %v", dErr)
		}
	} else {
		for _, result := range imageResults {
			err := searcher.Store.CreateImageResult(ctx, db.CreateImageParams{
				QueryID: queryId,
				Url:     result.Url,
				Data:    result.Data,
			})
			if err != nil {
				fmt.Printf("Error In storing ImageResult with url:%v err: %v", result.Url, err)
			}
		}

		fmt.Printf("RESULT: query:%s, len:%d page:%d perPage:%d\n", query, len(imageResults), page, number)

		dErr := searcher.Store.UpdateQuery(context.Background(), db.UpdateQueryArgs{
			Id:     queryId,
			Status: model.StatusSuccess,
		})
		if dErr != nil {
			fmt.Printf("Error in Updaing Success Request %v", dErr)
		}
	}
}
