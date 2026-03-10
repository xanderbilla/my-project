package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

const maxRequestBodySize = 1_048_576 // 1MB

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, status int, msg string) error {
	resp := Response{
		Status:  status,
		Message: "request failed",
		Error:   msg,
	}

	return JSON(w, status, resp)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) error {
	resp := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	return JSON(w, status, resp)
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("body must contain only one JSON object")
	}

	return nil
}
