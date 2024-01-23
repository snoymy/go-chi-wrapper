package router

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
)

func NewRouter() *Router {
    r := &Router{Router: chi.NewRouter()}

    return r
}

type HandlerFunc func(ResponseWriter, *Request) error

type ErrorHandlerFunc func(error, ResponseWriter, *Request)

type Router struct {
    chi.Router
    ErrorHandlers []ErrorHandlerFunc
}

func handlerWrapper(f HandlerFunc, errorHandlers []ErrorHandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := f(ResponseWriter{ResponseWriter: w}, &Request{Request: *r})

        if err != nil {
            for _, errorHandler := range(errorHandlers) {
                errorHandler(err, ResponseWriter{ResponseWriter: w}, &Request{Request: *r})
            }
        }
    }
}

func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
    r.Router.Use(middlewares...)
}

func (r *Router) UseErrorHandler(errorHandlers ...ErrorHandlerFunc) {
    if len(r.Router.Routes()) > 0 {
		panic("all error handler must be defined before routes on a mux")
	}

    r.ErrorHandlers = append(r.ErrorHandlers, errorHandlers...)
}

func (r *Router) With(middlewares ...func(http.Handler) http.Handler) *Router {
    return &Router{Router: r.Router.With(middlewares...), ErrorHandlers: r.ErrorHandlers}
}

func (r *Router) WithErrorHandler(errorHandlers ...ErrorHandlerFunc) *Router {
    return &Router{Router: r.Router, ErrorHandlers: append(r.ErrorHandlers, errorHandlers...)}
}

func (r *Router) Group(fn func(r *Router)) *Router {
    nr := r.With()
    if fn != nil {
		fn(nr)
	}

    return nr
}

func (r *Router) Route(pattern string, fn func(r *Router)) *Router {
    if fn == nil {
		panic(fmt.Sprintf("attempting to Route() a nil subrouter on '%s'", pattern))
	}

	subRouter := NewRouter()
	fn(subRouter)
	r.Mount(pattern, subRouter)

	return subRouter
}

func (r *Router) Mount(pattern string, h http.Handler) {
    r.Router.Mount(pattern, h)
}

func (r *Router) Handle(pattern string, h http.Handler) {
    r.Router.Handle(pattern, h)
}

func (r *Router) HandleFunc(pattern string, h HandlerFunc) {
    r.Router.HandleFunc(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Method(method, pattern string, h http.Handler) {
    r.Router.Method(method, pattern, h)
}
	
func (r *Router) MethodFunc(method, pattern string, h HandlerFunc) {
    r.Router.MethodFunc(method, pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Connect(pattern string, h HandlerFunc) {
    r.Router.Connect(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Delete(pattern string, h HandlerFunc) {
    r.Router.Delete(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Get(pattern string, h HandlerFunc) {
    r.Router.Get(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Head(pattern string, h HandlerFunc) {
    r.Router.Head(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Options(pattern string, h HandlerFunc) {
    r.Router.Options(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Patch(pattern string, h HandlerFunc) {
    r.Router.Patch(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Post(pattern string, h HandlerFunc) {
    r.Router.Post(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Put(pattern string, h HandlerFunc) {
    r.Router.Put(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) Trace(pattern string, h HandlerFunc) {
    r.Router.Trace(pattern, handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) NotFound(h HandlerFunc) {
    r.Router.NotFound(handlerWrapper(h, r.ErrorHandlers))
}

func (r *Router) MethodNotAllowed(h HandlerFunc) {
    r.Router.MethodNotAllowed(handlerWrapper(h, r.ErrorHandlers))
}
