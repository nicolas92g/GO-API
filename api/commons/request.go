package commons

import (
	"io"
	"net/http"
	"net/url"
	"projet/api/commons/apiError"
	"strings"
)

type RequestApi struct {
	r http.Request
}

func (r *RequestApi) Create(req *http.Request) {
	r.r = *req
}

func (r *RequestApi) GetMethod() string {
	return r.r.Method
}

func (r *RequestApi) GetUrl() string {
	return r.r.URL.Path
}

func (r *RequestApi) GetParameters() url.Values {
	return r.r.URL.Query()
}

func (r *RequestApi) GetBody() ([]byte, apiError.ApiError) {
	ret, err := io.ReadAll(r.r.Body)
	if err != nil {
		return []byte{}, apiError.InternalError{Msg: "Error while reading the request body : " + err.Error()}
	}
	return ret, nil
}

func (r *RequestApi) GetPathAt(index uint8) (string, apiError.ApiError) {
	ar := strings.Split(r.r.URL.Path, "/")
	if int(index) > (len(ar) - 1) {
		return "", apiError.InternalError{Msg: "Error while reading request URL : Bad index"}
	}
	return ar[index], nil
}

func (r *RequestApi) GetHeader() http.Header {
	return r.r.Header
}
