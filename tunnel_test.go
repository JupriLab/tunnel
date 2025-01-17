package tunnel_test

import (
	"net/http"
	"testing"

	"github.com/JupriLab/tunnel"
)

func TestRoute(t *testing.T) {
	route := tunnel.NewRoute()
	route.Get("/test", func(ctx *tunnel.ResponseWriter) error {
		return ctx.ResponseJson(tunnel.TResponseJson{
			Message: "running test",
		}, http.StatusOK)
	})

	// if err := route.Start(":3000"); err != nil {
	// 	t.Errorf("Test route error: %s\n", err.Error())
	// }

}
