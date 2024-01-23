package router

import (
    "encoding/json"
    "net/http"
)

type ResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (w *ResponseWriter) SetStatus(statusCode int) {
    w.statusCode = statusCode
}

func (w *ResponseWriter) Write(b []byte) error {
    if w.statusCode == 0 {
        w.statusCode = 200
    }

    w.WriteHeader(w.statusCode) // This must call after "w.Header().Add(....)"
    _, err := w.ResponseWriter.Write(b)
    return err
}

func (w *ResponseWriter) WriteString(s string) error {
    err := w.Write([]byte(s))
    return err
}

func (w *ResponseWriter) WriteAsJson(a any) error {
    var err error

    var j []byte
    j, err = json.Marshal(a)
    if err != nil {
        return err
    }

    w.Header().Add("Content-Type", "application/json")
    err = w.Write([]byte(string(j)))

    return err
}

