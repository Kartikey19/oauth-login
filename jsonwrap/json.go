package jsonwrap

import (
	"encoding/json"
	"net/http"
)

// MakeJSONResponse : Method to wrap a interface in a json byte response
func MakeJSONResponse(message string, payload interface{}, isSuccess bool) ([]byte, error) {
	data := make(map[string]interface{})
	data["message"] = message
	data["payload"] = payload
	data["success"] = 1
	if !isSuccess {
		data["success"] = 0
	}
	return json.Marshal(data)
}

// SendJSONHttpResponse : Http response writer helper
func SendJSONHttpResponse(w http.ResponseWriter, body []byte, httpStatus int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(body)
}
