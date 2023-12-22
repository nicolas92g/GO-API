package apiError

import "net/http"

type BddError struct {
	Msg string
}

func (ie BddError) Error() string {
	return ie.Msg
}

func (ie BddError) Code() int {
	return http.StatusInternalServerError
}
