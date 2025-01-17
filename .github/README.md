<p align="center">
    <em><b>Tunnel</b></em> is a web framework inspired by <a href="https://github.com/gofiber/fiber">Gofiber</a> and build on top of <a href="https://github.com/valyala/fasthttp">FastHttp</a>. We created this with the aim to isolate less used functionality when we use popular frameworks
</p>

---

## 🧪 **This is still in beta version!**

We're actively developing and improving the project, and we'd love your input! If you have ideas for features or improvements you'd like to see in **_Tunnel_**, feel free to create an [issue](https://github.com/JupriLab/tunnel/issues). Your feedback helps us shape the future of the project, so don't hesitate to share your thoughts! 😊

---

## 🛠️ Installation

Tunnel requires **Go version `1.23` or higher**

```bash
go mod init github.com/username/repo
```

after finish with setup your project, you can install Tunnel with command:

```bash
go get github.com/JupriLab/tunnel
```

## ✨ Quickstart

Getting started with **_Tunnel_** is super easy! Here’s a quick example to create a simple web server that sends back a JSON response with "Hello, World" in the message. 😊

```go title="Example"
package main

import (
	"net/http"

	"github.com/JupriLab/tunnel"
)

func main() {

	routes := tunnel.NewRoute()

	routes.Get("/", func(ctx *tunnel.ResponseWriter) error {
		return ctx.ResponseJson(tunnel.TResponseJson{
			Message: "Hello The World",
		}, http.StatusOK)
	})
}
```

## Examples

Here’s another example! If you want to organize your web server into groups or use prefixes for routes, this will help you get started easily.

```go title="Example 2"
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
```