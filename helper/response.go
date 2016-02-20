package helper

import (
	"strings"
	"github.com/qasico/beego/validation"
)

type Response struct {
	Format map[string]interface{}
}

var (
	Respond = Response{}
)

func (r *Response) SuccessWithData(total int64, data ...[]interface{}) {
	r.Format = make(map[string]interface{})
	r.Format["status"] = "success"
	r.Format["total"] = total

	if data != nil {
		r.Format["data"] = data[0]
	} else {
		r.Format["data"] = nil
	}
}

func (r *Response) SuccessWithModel(httpMethod string, model interface{}) {
	r.Format = make(map[string]interface{})
	r.Format["status"] = "success"

	if httpMethod == "POST" {
		r.Format["data"] = model
	} else if httpMethod == "GET" && model != nil {
		r.Format["data"] = model
	}
}

func (r *Response) Fail(errorData interface{}) {
	r.Format = make(map[string]interface{})
	r.Format["status"] = "failed"

	switch errorData.(type) {
	case string:
		r.Format["message"] = map[string]string{
			"error": ClearErrorPrefix(errorData.(string)),
		}
	default:
		r.Format["message"] = errorData
	}
}

func (r *Response) Validator(model interface{}) (bool, map[string]interface{}) {
	errorData := make(map[string]string)
	validator := validation.Validation{}

	passed, _ := validator.Valid(model)
	if !passed {
		for _, err := range validator.Errors {
			field := strings.Split(err.Key, ".")
			errorData[field[0]] = err.Message
		}

		r.Fail(errorData)
		return false, r.Format
	}

	return true, r.Format
}

func ClearErrorPrefix(s string) string {
	strToRemove := "<QuerySeter> "
	s = strings.TrimPrefix(s, strToRemove)
	return s
}