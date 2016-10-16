package utils

import (
	"encoding/json"
	"net/http"
)

type StringResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	OK      bool   `json:"ok"`
}

type DataResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	OK    bool        `json:"ok"`
}

func respondJSON(w http.ResponseWriter, i interface{}, code int) {
	// Shamelessly stolen from net/http.Error source code
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, err := json.Marshal(i)
	if err != nil {
		// nothing we can do here - http.Error just ignores it - should we do the same?
		panic(err)
	}

	_, err = w.Write(b)
	if err != nil {
		// nothing we can do here - http.Error just ignores it - should we do the same?
		panic(err)
	}
}

func RespondError(w http.ResponseWriter, errStr string, code int) {
	resp := StringResponse{
		Error:   errStr,
		Message: errStr,
		Code:    code,
		OK:      false,
	}

	respondJSON(w, resp, code)
}

func RespondSuccess(w http.ResponseWriter, msg string) {
	resp := StringResponse{
		Message: msg,
		Code:    200,
		OK:      true,
	}

	respondJSON(w, resp, 200)
}

func RespondDataError(w http.ResponseWriter, errStr string, data interface{}, code int) {
	resp := DataResponse{
		Error: errStr,
		Data:  data,
		Code:  code,
		OK:    false,
	}

	respondJSON(w, resp, code)
}

func RespondDataSuccess(w http.ResponseWriter, data interface{}) {
	respondJSON(w, data, 200)
}
