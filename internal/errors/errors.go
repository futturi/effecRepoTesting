package errors

import "errors"

var (
	ErrSongNotFound      = errors.New("song not found")
	ErrAnotherStatucCode = errors.New("error with sending request")
	ErrIncorrectRequest  = errors.New("incorrect request")
)

type ErrorMessage struct {
	Error string `json:"error"`
}
