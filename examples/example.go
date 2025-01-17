package main

import (
	"net/http"

	"github.com/JupriLab/tunnel"
)

func main() {

	routes := tunnel.NewRoute()

	routes.Get("/home", func(ctx *tunnel.ResponseWriter) error {
		return ctx.ResponseJson(tunnel.TResponseJson{
			Message: "Home Page",
		}, http.StatusOK)
	})

	auth := routes.Group("/auth")
	auth.Post("/login", func(ctx *tunnel.ResponseWriter) error {
		var body struct {
			Email    string
			Password string
		}

		if err := ctx.Bind(&body); err != nil {
			return ctx.ResponseJson(tunnel.TResponseJson{
				Message: "Failed Login",
				Data: map[string]string{
					"errorMessage": err.Error(),
				},
			}, http.StatusInternalServerError)
		}

		return ctx.ResponseJson(tunnel.TResponseJson{
			Message: "Success Login",
		}, http.StatusOK)
	})

	user := routes.Group("/user")
	user.Get("/", func(ctx *tunnel.ResponseWriter) error {
		customStatusCode := http.StatusAccepted

		return ctx.ResponseJson(tunnel.TResponseJson{
			Status:  &customStatusCode,
			Message: "Successfully get all users",
			Data: []map[string]string{
				{
					"id":    "1238012-123123",
					"email": "test@mail.com",
				},
			},
		}, http.StatusOK)
	})

	user.Get("/{id}", func(ctx *tunnel.ResponseWriter) error {
		userId := ctx.GetParam("id", nil)

		return ctx.ResponseJson(tunnel.TResponseJson{
			Message: "Successfully get user",
			Data: map[string]string{
				"id":    userId.(string),
				"email": "test@mail.com",
			},
		}, http.StatusOK)
	})
}
