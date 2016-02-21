package helper

import (
	"strings"

	"github.com/qasico/beego/validation"
)

type APIResponse struct {
	Code    int
	Data    interface{}
	Total   int64
	Status  string
	Message interface{}
}

func (APIRes *APIResponse) Success(total int64, data interface{}) {
	APIRes.Code = 200
	APIRes.Status = "success"
	APIRes.Total = total
	APIRes.Data = data
}

func (APIRes *APIResponse) Failed(code int, message interface{}) {
	APIRes.Code = code
	APIRes.Status = "failed"

	switch message.(type) {
	case string:
		APIRes.Message = map[string]string{
			"error": ClearErrorPrefix(message.(string)),
		}
	default:
		APIRes.Message = message
	}
}

func (APIRes *APIResponse) Validator(model interface{}) (bool) {
	errorData := make(map[string]string)
	validator := validation.Validation{}

	passed, _ := validator.Valid(model)
	if !passed {
		for _, err := range validator.Errors {
			field := strings.Split(err.Key, ".")
			errorData[field[0]] = err.Message
		}

		APIRes.Failed(304, errorData)
		return false
	}

	return true
}

func (APIRes *APIResponse) GetResponse(httpMethod string) (response map[string]interface{}) {
	response = make(map[string]interface{})
	response["status"] = APIRes.Status

	if APIRes.Status == "success" {
		if httpMethod == "GET" {
			response["data"] = APIRes.Data
			response["total"] = APIRes.Total
		} else {
			if APIRes.Data != nil {
				response["data"] = APIRes.Data
			}
		}
	} else {
		response["message"] = APIRes.Message
	}

	return response
}

func ClearErrorPrefix(s string) string {
	strToRemove := "<QuerySeter> "
	s = strings.TrimPrefix(s, strToRemove)
	return s
}