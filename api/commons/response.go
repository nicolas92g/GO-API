package commons

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projet/api/commons/apiError"
)

type ResponseApi struct {
	w http.ResponseWriter
}

func (r *ResponseApi) Create(res http.ResponseWriter) {
	r.w = res
	r.w.Header().Set("Content-Type", "application/json")
	r.w.Header().Set("Access-Control-Allow-Origin", "*")
	r.w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	r.w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (r *ResponseApi) SendMessage(msg string, code int) {
	r.w.WriteHeader(code)
	_, _ = fmt.Fprint(r.w, "{ \"message\": \""+msg+"\"}")
}

func (r *ResponseApi) SendError(err apiError.ApiError) bool {
	if err != nil {
		r.SendMessage(err.Error(), err.Code())
		return true
	}
	return false
}

func (r *ResponseApi) SendContent(content any, code int) {
	jsonContent, err := json.Marshal(content)

	if err != nil {
		r.SendMessage("Internal Server error : failed to Marshal json response", http.StatusInternalServerError)
		return
	}

	r.w.WriteHeader(code)
	_, _ = fmt.Fprint(r.w, string(jsonContent))
}
