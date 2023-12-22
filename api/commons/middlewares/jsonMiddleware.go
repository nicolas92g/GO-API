package middlewares

import (
	"encoding/json"
	"fmt"
	"projet/api/commons"
	apiError2 "projet/api/commons/apiError"
)

func JsonMiddleware(r *commons.RequestApi, body any) apiError2.ApiError {
	//test only for POST and PATCH requests
	if r.GetMethod() != "POST" && r.GetMethod() != "PATCH" {
		return nil
	}

	// check request body is in json
	if r.GetHeader().Get("Content-Type") != "application/json" {
		return apiError2.BadRequest{Msg: "Request Body has to be json"}
	}

	bodyBytes, err := r.GetBody()
	if err != nil {
		return apiError2.InternalError{Msg: "Failed to read request body"}
	}

	err2 := json.Unmarshal(bodyBytes, body)
	if err2 != nil {
		fmt.Print(string(bodyBytes))
		return apiError2.BadRequest{Msg: "JSON body can not be read : " + err.Error()}
	}
	return nil
}
