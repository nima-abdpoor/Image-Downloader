package db

import (
	"GoogleImageDownloader/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Driver                 string        `koanf:"driver"`
	Username               string        `koanf:"user"`
	Password               string        `koanf:"password"`
	Host                   string        `koanf:"host"`
	DBName                 string        `koanf:"name"`
	Port                   int           `koanf:"port"`
	ConnMaxLifeTimeMinutes int           `koanf:"maxlifetimeminutes"`
	MaxOpenCons            int           `koanf:"maxopencons"`
	MaxIdleCons            int           `koanf:"maxidlecons"`
	RetryConnection        time.Duration `koanf:"retry"`
}

type Store struct {
	q  *Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		q:  New(db),
		db: db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type QueryResult struct {
	Id      int64
	Title   string
	Status  string
	PerPage int32
	Page    int32
}

type UpdateQueryArgs struct {
	Id     int64
	Status string
}

type CreateImageParams struct {
	QueryID int64
	Url     string
	Data    string
}

func (store Store) CreateQuery(ctx context.Context, query string, perPage, page int32) (QueryResult, error) {
	var qResult QueryResult
	err := store.execTx(context.Background(), func(queries *Queries) error {
		var err error
		result, err := queries.CreateQuery(ctx, CreateQueryParams{
			Query:   query,
			Status:  model.StatusInProgress,
			PerPage: perPage,
			Page:    page,
		})
		if err != nil {
			return err
		}
		qResult.Id = result.ID
		return nil
	})
	return qResult, err
}

func (store Store) GetQueryByStatus(ctx context.Context, status string) ([]QueryResult, error) {
	var result []QueryResult
	if queries, err := store.q.GetQueryByStatus(ctx, status); err != nil {
		return nil, err
	} else {
		for _, q := range queries {
			result = append(result, QueryResult{
				Id:      q.ID,
				Title:   q.Query,
				Status:  q.Status,
				PerPage: q.PerPage,
				Page:    q.Page,
			})
		}
	}
	return result, nil
}

func (store Store) CreateImageResult(ctx context.Context, params CreateImageParams) error {
	err := store.execTx(context.Background(), func(queries *Queries) error {
		var err error
		_, err = queries.CreateImageResult(ctx, CreateImageResultParams{
			QueryID:   params.QueryID,
			ImageUrl:  sql.NullString{String: params.Url, Valid: true},
			ImageData: sql.NullString{String: params.Data, Valid: true},
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (store Store) UpdateQuery(ctx context.Context, arg UpdateQueryArgs) error {
	err := store.execTx(context.Background(), func(queries *Queries) error {
		var err error
		_, err = queries.UpdateQuery(ctx, UpdateQueryParams{
			ID:        arg.Id,
			Status:    arg.Status,
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
