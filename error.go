// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beego

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/qasico/beego/context"
	"github.com/qasico/beego/utils"
	"github.com/qasico/beego/helper"
)

const (
	errorTypeHandler = iota
	errorTypeController
)

// render default application error page with error and stack string.
func showErr(err interface{}, ctx *context.Context, stack string) {
	data := map[string]string{
		"error":          fmt.Sprintf("%v", err),
		"request_method": ctx.Input.Method(),
		"request_url":    ctx.Input.URI(),
	}

	renderError(500, ctx, data)
}

type errorInfo struct {
	controllerType reflect.Type
	handler        http.HandlerFunc
	method         string
	errorType      int
}

// ErrorMaps holds map of http handlers for each error string.
// there is 10 kinds default error(40x and 50x)
var ErrorMaps = make(map[string]*errorInfo, 10)

// ErrorHandler registers http.HandlerFunc to each http err code string.
// usage:
// 	beego.ErrorHandler("404",NotFound)
//	beego.ErrorHandler("500",InternalServerError)
func ErrorHandler(code string, h http.HandlerFunc) *App {
	ErrorMaps[code] = &errorInfo{
		errorType: errorTypeHandler,
		handler:   h,
		method:    code,
	}
	return BeeApp
}

// ErrorController registers ControllerInterface to each http err code string.
// usage:
// 	beego.ErrorController(&controllers.ErrorController{})
func ErrorController(c ControllerInterface) *App {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	for i := 0; i < rt.NumMethod(); i++ {
		methodName := rt.Method(i).Name
		if !utils.InSlice(methodName, exceptMethod) && strings.HasPrefix(methodName, "Error") {
			errName := strings.TrimPrefix(methodName, "Error")
			ErrorMaps[errName] = &errorInfo{
				errorType:      errorTypeController,
				controllerType: ct,
				method:         methodName,
			}
		}
	}
	return BeeApp
}

// show error string as simple text message.
// if error string is empty, show 503 or 500 error as default.
func exception(errCode string, ctx *context.Context) {
	atoi := func(code string) int {
		v, err := strconv.Atoi(code)
		if err == nil {
			return v
		}
		return 503
	}

	for _, ec := range []string{errCode, "503", "500"} {
		if h, ok := ErrorMaps[ec]; ok {
			executeError(h, ctx, atoi(ec))
			return
		}
	}

	renderError(atoi(errCode), ctx, http.StatusText(atoi(errCode)))
}

func executeError(err *errorInfo, ctx *context.Context, code int) {
	if err.errorType == errorTypeHandler {
		renderError(code, ctx, http.StatusText(code))
		return
	}

	if err.errorType == errorTypeController {
		ctx.Output.SetStatus(code)
		//Invoke the request handler
		vc := reflect.New(err.controllerType)
		execController, ok := vc.Interface().(ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface")
		}
		//call the controller init function
		execController.Init(ctx, err.controllerType.Name(), err.method, vc.Interface())

		//call prepare function
		execController.Prepare()

		execController.URLMapping()

		method := vc.MethodByName(err.method)
		method.Call([]reflect.Value{})

		//render template
		if BConfig.WebConfig.AutoRender {
			if err := execController.Render(); err != nil {
				panic(err)
			}
		}

		// finish all runrouter. release resource
		execController.Finish()
	}
}

func renderError(code int, ctx *context.Context, message ...interface{}) {
	response := helper.APIResponse{}
	response.Failed(code, message)

	ctx.Output.ContentType("json")
	ctx.ResponseWriter.WriteHeader(response.Code)
	ctx.Output.JSON(response.GetResponse("GET"), true, true)
}