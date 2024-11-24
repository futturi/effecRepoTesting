package store

import (
	"context"
	"effectiveMobile/internal/entities"
	"effectiveMobile/internal/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	Songs
}

func NewStore(db *gorm.DB) Store {
	return Store{
		Songs: NewStoreSongs(db),
	}
}

func InitDB(ctx context.Context, connString string) (*gorm.DB, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debugw("connecting to database")
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Errorw("error with opening database conn", zap.Error(err))
		return nil, err
	}

	if err := db.AutoMigrate(&entities.Song{}); err != nil {
		log.Errorw("error with migrating database scheme", zap.Error(err))
		return nil, err
	}

	log.Debug("database is connected")

	return db, nil
}

func ShutdownDb(ctx context.Context, db *gorm.DB) error {
	log := logger.LoggerFromContext(ctx)
	log.Debugw("closing database conn")
	dbSQL, err := db.DB()
	if err != nil {
		log.Errorw("error with closing database conn", zap.Error(err))
		return err
	}
	log.Debugw("database conn is closed")
	return dbSQL.Close()
}

type Songs interface {
	InsertSong(ctx context.Context, req entities.Song) (int, error)
	GetSongs(ctx context.Context, filters entities.Song, limit, offset int) ([]entities.Song, error)
	DeleteSong(ctx context.Context, songId int) error
	GetTextSong(ctx context.Context, songId int) (string, error)
	UpdateSong(ctx context.Context, songId int, song entities.SongUpdate) error
}
