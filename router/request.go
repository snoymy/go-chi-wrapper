package router 

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

type Request struct {
    http.Request
}

func (r *Request) URLParam(param string) string {
    return chi.URLParam(&r.Request, param)
}

func (r *Request) URLQuery(query string) string {
    return r.URL.Query().Get(query)
}

func (r *Request) URLQueries(query string) []string {
    return r.URL.Query()[query]
}

