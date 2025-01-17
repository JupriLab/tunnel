package tunnel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

func NewRoute() *Route {
	return &Route{
		prefix:     "",
		routes:     make(map[string]map[string]func(*fasthttp.RequestCtx)),
		middleware: []func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx){},
	}
}

func (route *Route) Group(prefix string) *Route {
	return &Route{
		prefix:     route.prefix + strings.TrimRight(prefix, "/"),
		routes:     route.routes,
		middleware: append([]func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx){}, route.middleware...),
	}
}

func (r *Route) Middlewares(middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	// append middleware to the route
	r.middleware = append(r.middleware, middlewares...)
}

func (r *Route) Handle(method string, path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	fullPath := r.prefix + path

	// Initialize the route map for the full path if not present
	if r.routes[fullPath] == nil {
		r.routes[fullPath] = make(map[string]func(*fasthttp.RequestCtx))
	}

	// Check if the route already exists for the method
	if _, exists := r.routes[fullPath][method]; exists {
		panic(fmt.Sprintf("Route already registered: %s %s", method, fullPath))
	}

	finalHandler := func(ctx *fasthttp.RequestCtx) {
		handler(&ResponseWriter{ctx})
	}

	for _, middleware := range r.middleware {
		finalHandler = middleware(finalHandler)
	}

	for _, middleware := range middlewares {
		finalHandler = middleware(finalHandler)
	}

	r.routes[fullPath][method] = finalHandler
}

// Helper methods to define HTTP methods for routes
func (r *Route) Get(path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	r.Handle("GET", path, handler, middlewares...)
}

func (r *Route) Post(path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	r.Handle("POST", path, handler, middlewares...)
}

func (r *Route) Put(path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	r.Handle("PUT", path, handler, middlewares...)
}

func (r *Route) Option(path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	r.Handle("OPTION", path, handler, middlewares...)
}

func (r *Route) Delete(path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	r.Handle("DELETE", path, handler, middlewares...)
}

func (r *Route) CustomMethod(method string, path string, handler func(ctx *ResponseWriter) error, middlewares ...func(func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx)) {
	r.Handle(strings.ToUpper(method), path, handler, middlewares...)
}

func (r *Route) routeRegister(ctx *fasthttp.RequestCtx) {
	routePath := string(ctx.Path())
	method := string(ctx.Method())

	for routePattern, methods := range r.routes {
		// Handle static routes (without dynamic parameters)
		if routePattern == routePath {
			if routeHandle, exists := methods[method]; exists {
				routeHandle(ctx)
				return
			}
		}

		// Handle dynamic routes (with parameters)
		if matched, _ := regexp.MatchString("^"+regexp.MustCompile("{([a-zA-Z0-9_-]+)}").ReplaceAllString(routePattern, "([a-zA-Z0-9_-]+)")+"$", routePath); matched {
			// Convert route pattern to regular expression for dynamic path matching
			regPath := regexp.MustCompile("^" + regexp.MustCompile("{([a-zA-Z0-9_-]+)}").ReplaceAllString(routePattern, "([a-zA-Z0-9_-]+)") + "$")

			// Extract parameter names from the route pattern
			namesRegex := regexp.MustCompile("{([a-zA-Z0-9_-]+)}").FindAllStringSubmatch(routePattern, -1)

			// Create a map to hold the extracted parameters
			params := make(map[string]string)

			// Get the dynamic values from the matched path
			valueRegex := regPath.FindStringSubmatch(routePath)[1:]

			// Map the extracted values to parameter names
			for i, value := range valueRegex {
				name := namesRegex[i][1]
				params[name] = value
			}

			// Set the parameters in the context for later use
			ctx.SetUserValue("params", params)

			// If a handler exists for the current HTTP method, invoke it
			if routeHandle, exists := methods[method]; exists {
				routeHandle(ctx)
				return
			}
		}
	}

	ctxUser := &ResponseWriter{ctx}
	ctxUser.ResponseNotFound()
	return
}

func (r *Route) Start(addr string) error {
	if r.routes == nil {
		panic("Routes is not initialized")
	}

	return fasthttp.ListenAndServe(addr, r.routeRegister)
}
