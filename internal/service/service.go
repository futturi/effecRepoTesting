package service

import (
	"context"
	"effectiveMobile/internal/entities"
	"effectiveMobile/internal/store"
)

type Service struct {
	Songs
}

func NewService(store *store.Store, songApiUrl string, apiTimeout int) *Service {
	return &Service{
		Songs: NewSongService(store.Songs, songApiUrl, apiTimeout),
	}
}

type Songs interface {
	InsertSong(ctx context.Context, req entities.SongRequest) (int, error)
	GetSongs(ctx context.Context, filters entities.Song, limit, offset string) ([]entities.Song, error)
	DeleteSong(ctx context.Context, songId string) error
	GetTextSong(ctx context.Context, lineInVerse, page, limit, songId string) (string, error)
	UpdateSong(ctx context.Context, songId string, song entities.SongUpdate) error
}
