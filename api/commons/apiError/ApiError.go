package apiError

type ApiError interface {
	Error() string
	Code() int
}
