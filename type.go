package tunnel

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

type Route struct {
	prefix     string
	routes     map[string]map[string]func(*fasthttp.RequestCtx) // map to store routes for specific HTTP methods
	middleware []func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)
}

type TResponseJson struct {
	Message string `json:"message"`
	Status  *int   `json:"status"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type ResponseWriter struct {
	*fasthttp.RequestCtx
}

func (ctx *ResponseWriter) Bind(body interface{}) error {
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		return err
	}

	return nil
}

func (ctx *ResponseWriter) GetParam(key string, defaultValue any) any {
	params := ctx.UserValue("params").(map[string]string)

	if param, exist := params[key]; exist {
		return param
	}

	return defaultValue
}

func (ctx *ResponseWriter) ResponseJson(data TResponseJson, statusCode int) error {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)

	if data.Status == nil {
		data.Status = &statusCode
	}

	err := json.NewEncoder(ctx).Encode(data)

	if err != nil {
		ctx.Error("Failed to encode response: ", http.StatusInternalServerError)
	}

	return err
}

func (ctx *ResponseWriter) ResponseNotFound() {
	ctx.ResponseJson(TResponseJson{
		Data:    nil,
		Message: "Not found",
	}, http.StatusNotFound)
}
