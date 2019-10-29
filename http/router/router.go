// Package router provides utils to add middleware handlers to
// http handler functions.
package router

import (
	"net/http"

	"github.com/pkg/errors"
	mw "github.com/trencat/goutils/http/middleware"
)

type handler struct {
	middleware []mw.Middleware
	Func       http.HandlerFunc
}

type Router map[string]handler

// Add adds middleware handlers to the given route.
// Middleware is executed in the order provided.
func (r *Router) Add(route string, mw ...mw.Middleware) {
	handler := (*r)[route]
	handler.middleware = append(handler.middleware, mw...)
	(*r)[route] = handler // TODO Check Do we need this line?
}

// HandleFunc sets the handler function to the given route.
// There can be only one handler function per route.
func (r *Router) HandleFunc(route string, h http.HandlerFunc) {
	handler := (*r)[route]
	handler.Func = h
	(*r)[route] = handler // TODO Check Do we need this line?
}

// Build registers all middleware and http handlers of the router
// to a multiplexer. If nil, the DefaultServeMux is used.
// This method must be called when all middleware and http handlers
// have been added to the router.
func (r *Router) Build(m *http.ServeMux) error {
	for route, handler := range *r {
		handlerFunc := handler.Func

		if handlerFunc == nil {
			return errors.Errorf("No HandlerFunc specified for route %s", route)
		}

		if len(handler.middleware) >= 1 {
			handlerFunc = build(handler.middleware, handlerFunc)
		}

		if m == nil {
			http.HandleFunc(route, handlerFunc)
		} else {
			m.HandleFunc(route, handlerFunc)
		}

	}

	return nil
}

func build(mw []mw.Middleware, handler http.HandlerFunc) http.HandlerFunc {
	if len(mw) == 1 {
		return mw[0](handler)
	}
	return mw[0](build(mw[1:], handler))
}
