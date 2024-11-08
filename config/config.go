package config

import (
	db "GoogleImageDownloader/db/sqlc"
	"GoogleImageDownloader/repository"
	"GoogleImageDownloader/scheduler"
	"GoogleImageDownloader/service"
)

type Config struct {
	PostgreSQL            db.Config         `koanf:"postgresql"`
	ReliableImageSearcher service.Config    `koanf:"image_searcher"`
	Scheduler             scheduler.Config  `koanf:"scheduler"`
	Pixel                 repository.Config `koanf:"pixel"`
}
