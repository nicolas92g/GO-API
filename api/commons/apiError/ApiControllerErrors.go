package apiError

import "net/http"

type InternalError struct {
	Msg string
}

func (ie InternalError) Error() string {
	return ie.Msg
}

func (ie InternalError) Code() int {
	return http.StatusInternalServerError
}

type BadRequest struct {
	Msg string
}

func (ie BadRequest) Error() string {
	return ie.Msg
}

func (ie BadRequest) Code() int {
	return http.StatusBadRequest
}

type BadURL struct {
	Msg string
}

func (ie BadURL) Error() string {
	return ie.Msg
}
func (ie BadURL) Code() int {
	return http.StatusBadRequest
}

type BadMethod struct {
	Msg string
}

func (ie BadMethod) Error() string {
	return ie.Msg
}
func (ie BadMethod) Code() int {
	return http.StatusBadRequest
}
