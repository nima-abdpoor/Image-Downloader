package main

import (
	"GoogleImageDownloader/config"
	db "GoogleImageDownloader/db/sqlc"
	"GoogleImageDownloader/repository"
	"GoogleImageDownloader/scheduler"
	"GoogleImageDownloader/service"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

func main() {
	cfg := config.Load("config.yml")

	fmt.Println(cfg)

	conn, err := getPostgreSQLConnection(cfg.PostgreSQL)
	if err != nil {
		log.Fatalf("Failed to Connect to Database %v", err)
	}

	store := db.NewStore(conn)

	searchEngine := repository.New(repository.Config{
		ApiKey: cfg.Pixel.ApiKey,
		File:   cfg.Pixel.File,
	})

	imageSearcherSvc := service.New(store, cfg.ReliableImageSearcher)

	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sch := scheduler.New(imageSearcherSvc, searchEngine, store, cfg.Scheduler)
		wg.Add(1)
		sch.Start(done, &wg)
	}()

	for true {

		fmt.Println("Enter Your Query...")
		query := "cat"
		fmt.Scan(&query)

		fmt.Println("Enter Max Number...")
		number := 80
		fmt.Scan(&number)

		page := number / cfg.Pixel.Min
		for i := 1; i <= page+1; i++ {
			if i == page+1 {
				go imageSearcherSvc.SearchImage(searchEngine, query, int32(number-page*cfg.Pixel.Min), int32(i), &wg)
			} else {
				go imageSearcherSvc.SearchImage(searchEngine, query, int32(cfg.Pixel.Min), int32(i), &wg)
			}
		}
	}

	wg.Wait()

	// todo add stateful shutdown
}

func getPostgreSQLConnection(config db.Config) (*sql.DB, error) {
	dbSource := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)
	if conn, err := sql.Open(config.Driver, dbSource); err != nil {
		//todo retry policy should be with some delay
		fmt.Println("Error in connecting to postgresql", err)
		return nil, err
		//return getPostgreSQLConnection(config)
	} else {
		return conn, nil
	}
}
