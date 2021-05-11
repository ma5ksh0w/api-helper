package helper

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code     int    `json:"code,omitempty"`
	Message  string `json:"message"`
	httpCode int
}

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

func NewError(code, httpCode int, msg string) *Response {
	return &Response{
		Error: &Error{
			httpCode: httpCode,
			Code:     code,
			Message:  msg,
		},
	}
}

func WriteError(rw http.ResponseWriter, code, httpCode int, msg string) error {
	return NewError(code, httpCode, msg).WriteTo(rw)
}

func NewOK(result interface{}) *Response {
	return &Response{
		Success: true,
		Result:  result,
	}
}

func WriteOK(rw http.ResponseWriter, result interface{}) error {
	return NewOK(result).WriteTo(rw)
}

func (r *Response) WriteTo(rw http.ResponseWriter) error {
	rw.Header().Set("Content-type", "application/json")
	data, err := json.Marshal(r)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error":{"code":-255, "message": "cannot marshal response"}}`))
		return err
	}

	if r.Error != nil {
		if r.Error.httpCode == 0 {
			r.Error.httpCode = http.StatusInternalServerError
		}

		rw.WriteHeader(r.Error.httpCode)
		rw.Write(data)
		return nil
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
	return nil
}
