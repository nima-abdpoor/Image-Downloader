package scheduler

import (
	db "GoogleImageDownloader/db/sqlc"
	"GoogleImageDownloader/model"
	"GoogleImageDownloader/repository"
	"GoogleImageDownloader/service"
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	sch        *gocron.Scheduler
	searchSvc  service.SearchImage
	repository repository.SearchEngine
	store      *db.Store
	conf       Config
}

type Config struct {
	Interval     int           `koanf:"interval"`
	RetryTimeout time.Duration `koanf:"retry_time_out"`
}

func New(searchEng service.SearchImage, repository repository.SearchEngine, store *db.Store, conf Config) Scheduler {
	return Scheduler{
		sch:        gocron.NewScheduler(time.UTC),
		searchSvc:  searchEng,
		repository: repository,
		store:      store,
		conf:       conf,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	log.Println("starting scheduler...")
	defer wg.Done()

	if _, err := s.sch.Every(s.conf.Interval).Second().Do(s.RetryFailedQueries); err != nil {
		log.Println("error in calling RetryFailedQueries...", err)
	}
	s.sch.StartAsync()

	<-done
	log.Printf("stopping scheduler...")
	s.sch.Stop()
}

func (s Scheduler) RetryFailedQueries() {
	ctx, cancel := context.WithTimeout(context.Background(), s.conf.RetryTimeout)
	defer cancel()

	results, err := s.store.GetQueryByStatus(ctx, model.StatusFailed)

	if err != nil {
		fmt.Println("Error in Finding Failed Queries!", err)
		return
	}

	var wg sync.WaitGroup
	for _, query := range results {
		// todo we should remove repetitive queries.
		s.searchSvc.RetrySearchImage(s.repository, query.Title, query.Id, query.PerPage, query.Page, &wg)
	}
}
