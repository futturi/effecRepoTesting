package service

import (
	"context"
	"effectiveMobile/internal/entities"
	"effectiveMobile/internal/errors"
	"effectiveMobile/internal/logger"
	"effectiveMobile/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SongService struct {
	store      store.Songs
	songApiUrl string
	httpClient http.Client
}

func NewSongService(store store.Songs, songApiUrl string, apiTimeout int) *SongService {
	client := http.Client{
		Timeout: time.Duration(apiTimeout) * time.Second,
	}
	return &SongService{
		store:      store,
		songApiUrl: songApiUrl,
		httpClient: client,
	}
}

func (s *SongService) InsertSong(ctx context.Context, insertReq entities.SongRequest) (int, error) {
	log := logger.LoggerFromContext(ctx)
	log = log.With("song", insertReq.Song, "group", insertReq.Group)
	requestUrl := fmt.Sprintf("%s/info?group=%s&song=%s", s.songApiUrl, insertReq.Group, insertReq.Song)
	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Errorw("error with creating request", zap.Error(err))
		return 0, err
	}
	log.Infow("doing request to song api", "requestUrl", requestUrl)
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Errorw("error with sending request", zap.Error(err))
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		log.Errorw("error with sending request", zap.Error(err))
		return 0, errors.ErrAnotherStatucCode
	}

	var songDetail entities.Song
	err = json.NewDecoder(response.Body).Decode(&songDetail)
	if err != nil {
		log.Errorw("error with decoding response", zap.Error(err))
		return 0, err
	}

	songDetail.GroupName = insertReq.Group
	songDetail.Song = insertReq.Song
	log.Info("song is received")
	id, err := s.store.InsertSong(ctx, songDetail)
	if err != nil {
		log.Errorw("error with inserting song", zap.Error(err))
		return 0, err
	}

	log.Info("song is inserted")
	return id, nil
}

func (s *SongService) GetSongs(ctx context.Context, filters entities.Song, limit, offset string) ([]entities.Song, error) {
	log := logger.LoggerFromContext(ctx)
	log = log.With("filters", filters, "limit", limit, "offset", offset)
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "1"
	}
	_, err := time.Parse("02.01.2006", filters.ReleaseDate)
	if err != nil {
		log.Errorw("error with parsing release date", zap.Error(err))
		return nil, errors.ErrIncorrectRequest
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Errorw("error with converting limit to int", zap.Error(err))
		return nil, errors.ErrIncorrectRequest
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Errorw("error with converting offset to int", zap.Error(err))
		return nil, errors.ErrIncorrectRequest
	}
	if limitInt == 0 {
		limitInt = 10
	}
	return s.store.GetSongs(ctx, filters, limitInt, offsetInt)
}

func (s *SongService) DeleteSong(ctx context.Context, songId string) error {
	log := logger.LoggerFromContext(ctx)
	log = log.With("songId", songId)
	songIdInt, err := strconv.Atoi(songId)

	if err != nil {
		log.Errorw("error with converting id to int", zap.Error(err))
		return errors.ErrIncorrectRequest
	}

	return s.store.DeleteSong(ctx, songIdInt)
}

func (s *SongService) GetTextSong(ctx context.Context, lineInVerse, page, limit, songId string) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log = log.With("songId", songId)
	songIdInt, err := strconv.Atoi(songId)
	if err != nil {
		log.Errorw("error with converting id to int", zap.Error(err))
		return "", err
	}
	if lineInVerse == "" {
		lineInVerse = "4"
	}
	if page == "" || page == "0" {
		page = "1"
	}
	if limit == "" || limit == "0" {
		limit = "10"
	}
	lineInVerseInt, err := strconv.Atoi(lineInVerse)
	if err != nil {
		log.Errorw("error with converting lineInVerse to int", zap.Error(err))
		return "", errors.ErrIncorrectRequest
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Errorw("error with converting page to int", zap.Error(err))
		return "", errors.ErrIncorrectRequest
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Errorw("error with converting limit to int", zap.Error(err))
		return "", errors.ErrIncorrectRequest
	}

	fullText, err := s.store.GetTextSong(ctx, songIdInt)
	if err != nil {
		log.Errorw("error with getting song text", zap.Error(err))
		return "", errors.ErrIncorrectRequest
	}
	textSplited := strings.Split(fullText, "\n")

	pageLength := limitInt * lineInVerseInt

	startIndex := (pageInt - 1) * pageLength
	endIndex := startIndex + pageLength

	if startIndex >= len(textSplited) {
		return "", nil
	}
	if endIndex > len(textSplited) {
		endIndex = len(textSplited)
	}

	pageText := textSplited[startIndex:endIndex]

	result := make([]string, 0)
	for i := 0; i < len(pageText); i += lineInVerseInt {
		verseEnd := i + lineInVerseInt
		if verseEnd > len(pageText) {
			verseEnd = len(pageText)
		}
		result = append(result, strings.Join(pageText[i:verseEnd], "\n"))
	}

	return strings.Join(result, "\n"), nil
}

func (s *SongService) UpdateSong(ctx context.Context, songId string, song entities.SongUpdate) error {
	log := logger.LoggerFromContext(ctx)
	log = log.With("songId", songId)
	_, err := time.Parse("02.01.2006", *song.ReleaseDate)
	if err != nil {
		log.Errorw("error with parsing release date", zap.Error(err))
		return errors.ErrIncorrectRequest
	}
	songIdInt, err := strconv.Atoi(songId)
	if err != nil {
		log.Errorw("error with converting id to int", zap.Error(err))
		return errors.ErrIncorrectRequest
	}
	return s.store.UpdateSong(ctx, songIdInt, song)
}
