package handler

import (
	"effectiveMobile/internal/entities"
	projectError "effectiveMobile/internal/errors"
	"effectiveMobile/internal/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary InsertSong
// @Tags songs
// @Description insert song
// @ID insert song
// @Accept json
// @Produce json
// @Param input body entities.SongRequest true "Song"
// @Success 200 {object} entities.InsertResponse
// @Failure 400 {object} errors.ErrorMessage
// @Failure 500 {object} errors.ErrorMessage
// @Router /api/insertsong [post]
func (h *Handler) InsertSong(c *gin.Context) {
	log := logger.LoggerFromContext(c)
	var req entities.SongRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
			Error: "error with binding json",
		})
		log.Errorw("error with binding json", zap.Error(err))
		return
	}

	id, err := h.service.InsertSong(c.Request.Context(), req)

	if err != nil {
		if errors.Is(err, projectError.ErrAnotherStatucCode) {
			c.JSON(http.StatusInternalServerError, projectError.ErrorMessage{
				Error: "error with sending request",
			})
			log.Errorw("error with sending request", zap.Error(err))
			return
		}
		c.JSON(http.StatusInternalServerError, projectError.ErrorMessage{
			Error: "error with inserting song",
		})
		log.Errorw("error with inserting song", zap.Error(err))
		return
	}
	log.Info("song is inserted")
	c.JSON(http.StatusOK, entities.InsertResponse{
		ID: id,
	})
}

// @Summary GetSongs
// @Tags songs
// @Description get songs
// @ID get songs
// @Accept json
// @Produce json
// @Param limit query int false "limit"
// @Param page query int false   "page"
// @Param groupName query string false "groupName"
// @Param song query string false "song"
// @Param releaseDate query string false "releaseDate"
// @Param text query string false "text"
// @Param link query string false "link"
// @Success 200 {object} []entities.Song
// @Failure 400 {object} errors.ErrorMessage
// @Failure 404 {object} errors.ErrorMessage
// @Failure 500 {object} errors.ErrorMessage
// @Router /api/getsongs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	log := logger.LoggerFromContext(c)
	var filters entities.Song

	limit := c.Query("limit")
	page := c.Query("page")
	groupName := c.Query("groupName")
	song := c.Query("song")
	releaseDate := c.Query("releaseDate")
	text := c.Query("text")
	link := c.Query("link")

	filters.GroupName = groupName
	filters.Song = song
	filters.ReleaseDate = releaseDate
	filters.Text = text
	filters.Link = link

	result, err := h.service.GetSongs(c.Request.Context(), filters, limit, page)

	if err != nil {
		if errors.Is(err, projectError.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, projectError.ErrorMessage{
				Error: "song not found",
			})
			log.Errorw("song not found", zap.Error(err))
			return
		} else if errors.Is(err, projectError.ErrIncorrectRequest) {
			c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
				Error: "incorrect request",
			})
			log.Errorw("incorrect request", zap.Error(err))
			return
		}
		c.JSON(http.StatusInternalServerError, projectError.ErrorMessage{
			Error: "error with getting songs",
		})
		log.Errorw("error with getting songs", zap.Error(err))
		return
	}
	log.Infow("songs are got")
	c.JSON(http.StatusOK, result)
}

// @Summary DeleteSong
// @Tags songs
// @Description delete song
// @ID delete song
// @Accept json
// @Produce json
// @Param id path int true "songId"
// @Success 200 {object} entities.DeleteResponse
// @Failure 400 {object} errors.ErrorMessage
// @Failure 404 {object} errors.ErrorMessage
// @Failure 500 {object} errors.ErrorMessage
// @Router /api/deletesong/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	log := logger.LoggerFromContext(c)
	songId := c.Param("id")
	if songId == "" {
		c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
			Error: "error with getting song id",
		})
		log.Errorw("error with getting song id")
		return
	}

	err := h.service.DeleteSong(c.Request.Context(), songId)

	if err != nil {
		if errors.Is(err, projectError.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, projectError.ErrorMessage{
				Error: "song not found",
			})
			log.Errorw("song not found", zap.Error(err))
			return
		}
		c.JSON(http.StatusInternalServerError, projectError.ErrorMessage{
			Error: "error with deleting song",
		})
		log.Errorw("error with deleting song", zap.Error(err))
		return
	}
	log.Infow("song is deleted")
	c.JSON(http.StatusOK, entities.DeleteResponse{
		Status: true,
	})
}

// @Summary GetTextSong
// @Tags songs
// @Description get text song
// @ID get text song
// @Accept json
// @Produce json
// @Param id path int true "songId"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param lineInVerse query int false "lineInVerse"
// @Success 200 {object} entities.TextResponse
// @Failure 400 {object} errors.ErrorMessage
// @Failure 404 {object} errors.ErrorMessage
// @Failure 500 {object} errors.ErrorMessage
// @Router /api/gettext/{id} [get]
func (h *Handler) GetTextSong(c *gin.Context) {
	log := logger.LoggerFromContext(c)
	songId := c.Param("id")
	if songId == "" {
		c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
			Error: "error with getting song id",
		})
		log.Errorw("error with getting song id")
		return
	}

	lineInVerse := c.Query("lineInVerse")
	page := c.Query("page")
	limit := c.Query("limit")
	text, err := h.service.GetTextSong(c.Request.Context(), lineInVerse, page, limit, songId)

	if err != nil {
		if errors.Is(err, projectError.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, projectError.ErrorMessage{
				Error: "song not found",
			})
			log.Errorw("song not found", zap.Error(err))
			return
		} else if errors.Is(err, projectError.ErrIncorrectRequest) {
			c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
				Error: "incorrect request",
			})
			log.Errorw("incorrect request", zap.Error(err))
			return
		}
		c.JSON(http.StatusInternalServerError, projectError.ErrorMessage{
			Error: "error with getting song",
		})
		log.Errorw("error with getting song", zap.Error(err))
		return
	}
	log.Infow("song is got")
	c.JSON(http.StatusOK, entities.TextResponse{
		Text: text,
	})
}

// @Summary UpdateSong
// @Tags songs
// @Description update song
// @ID update song
// @Accept json
// @Produce json
// @Param id path int true "songId"
// @Param input body entities.SongUpdate true "Song"
// @Success 200 {object} entities.TextResponse
// @Failure 400 {object} errors.ErrorMessage
// @Failure 404 {object} errors.ErrorMessage
// @Failure 500 {object} errors.ErrorMessage
// @Router /api/updatesong/{id} [patch]
func (h *Handler) UpdateSong(c *gin.Context) {
	log := logger.LoggerFromContext(c)
	songId := c.Param("id")
	if songId == "" {
		c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
			Error: "error with getting song id",
		})
		log.Errorw("error with getting song id")
		return
	}

	var song entities.SongUpdate

	if err := c.BindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
			Error: "error with binding json",
		})
		log.Errorw("error with binding json", zap.Error(err))
		return
	}

	err := h.service.UpdateSong(c.Request.Context(), songId, song)

	if err != nil {
		if errors.Is(err, projectError.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, projectError.ErrorMessage{
				Error: "song not found",
			})
			log.Errorw("song not found", zap.Error(err))
			return
		} else if errors.Is(err, projectError.ErrIncorrectRequest) {
			c.JSON(http.StatusBadRequest, projectError.ErrorMessage{
				Error: "incorrect request",
			})
			log.Errorw("incorrect request", zap.Error(err))
			return
		}
		c.JSON(http.StatusInternalServerError, projectError.ErrorMessage{
			Error: "error with updating song",
		})
		log.Errorw("error with updating song", zap.Error(err))
		return

	}
	log.Infow("song is updated")
	c.JSON(http.StatusOK, entities.UpdateResponse{
		ID: songId,
	})
}
