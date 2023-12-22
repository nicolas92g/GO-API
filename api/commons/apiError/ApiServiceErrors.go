package apiError

import "net/http"

type ValueNotFound struct {
	Msg string
}

func (ie ValueNotFound) Error() string {
	return ie.Msg
}

func (ie ValueNotFound) Code() int {
	return http.StatusBadRequest
}

type AlreadyExists struct {
	Msg string
}

func (ie AlreadyExists) Error() string {
	return ie.Msg
}

func (ie AlreadyExists) Code() int {
	return http.StatusBadRequest
}
