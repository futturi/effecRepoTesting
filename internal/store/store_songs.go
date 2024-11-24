package store

import (
	"context"
	"effectiveMobile/internal/entities"
	"effectiveMobile/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StoreSongs struct {
	db *gorm.DB
}

func NewStoreSongs(db *gorm.DB) *StoreSongs {
	return &StoreSongs{
		db: db,
	}
}

func (r *StoreSongs) InsertSong(ctx context.Context, req entities.Song) (int, error) {
	log := logger.LoggerFromContext(ctx)
	if err := r.db.Create(&req).Error; err != nil {
		log.Errorw("error with inserting song", zap.Error(err))
		return 0, err
	}
	log.Infow("song is inserted")
	return int(req.ID), nil
}

func (r *StoreSongs) GetSongs(ctx context.Context, filters entities.Song, limit, offset int) ([]entities.Song, error) {
	log := logger.LoggerFromContext(ctx)
	var songs []entities.Song
	query := r.db.Model(&entities.Song{})
	if filters.GroupName != "" {
		query = query.Where("group_name LIKE ?", "%"+filters.GroupName+"%")
	}
	if filters.Song != "" {
		query = query.Where("song LIKE ?", "%"+filters.Song+"%")
	}
	if filters.ReleaseDate != "" {
		query = query.Where("release_date LIKE ?", "%"+filters.ReleaseDate+"%")
	}
	if filters.Text != "" {
		query = query.Where("text LIKE ?", "%"+filters.Text+"%")
	}
	if filters.Link != "" {
		query = query.Where("text LIKE ?", "%"+filters.Link+"%")
	}

	query = query.Offset((offset - 1) * limit).Limit(limit)
	tx := query.Find(&songs)

	if tx.Error != nil {
		log.Errorw("error with getting songs", zap.Error(tx.Error))
		return nil, tx.Error
	}

	log.Infow("songs are got")
	return songs, nil
}

func (r *StoreSongs) DeleteSong(ctx context.Context, songId int) error {
	log := logger.LoggerFromContext(ctx)
	if err := r.db.Delete(&entities.Song{}, songId).Error; err != nil {
		log.Errorw("error with deleting song", zap.Error(err))
		return err
	}
	log.Infow("song is deleted", "songId", songId)
	return nil
}

func (r *StoreSongs) GetTextSong(ctx context.Context, songId int) (string, error) {
	log := logger.LoggerFromContext(ctx)
	var song entities.Song
	if err := r.db.Where("id = ?", songId).First(&song).Error; err != nil {
		log.Errorw("error with getting song", zap.Error(err))
		return "", err
	}
	log.Infow("song is got", "songId", songId, "text", song.Text)
	return song.Text, nil
}

func (r *StoreSongs) UpdateSong(ctx context.Context, songId int, song entities.SongUpdate) error {
	log := logger.LoggerFromContext(ctx)
	if err := r.db.Model(&entities.Song{}).Where("id = ?", songId).Updates(song).Error; err != nil {
		log.Errorw("error with updating song", zap.Error(err))
		return err
	}
	log.Infow("song is updated", "songId", songId)
	return nil
}
