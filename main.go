package main

import (
    "fmt"
    "net/http"
    "myapp/router"
    "errors"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := router.NewRouter()

    //r.UseErrorHandler(basicErrorHandler)

    r.Get("/", helloWorld)
    r.Get("/json", helloJson)
    r.WithErrorHandler(basicErrorHandler).Get("/basic-error-handler", helloError)
    r.WithErrorHandler(myErrorHandler).Get("/common-error-handler", helloError)
    r.WithErrorHandler(myErrorHandler).Get("/my-error-handler", helloMyError)

    r.Route("/error-handler", func(r *router.Router){
        r.UseErrorHandler(basicErrorHandler)
        r.Use(middleware.Logger)

        r.Get("/common-error", helloError)
        r.Get("/my-error", helloMyError)
    })

    fmt.Println("Starting server at port 3000")
    http.ListenAndServe(":3000", r)
}

type ResponseObject struct {
    Status Status           `json:"status"`
    Data any                `json:"data"`
}

type Status struct {
    Code int            `json:"code"`
    Description string  `json:"description"`
}

type MyError struct {
    HttpStatus int
    StatusCode int
    Description string
}

func (e MyError) Error() string {
    return e.Description
}

func basicErrorHandler(e error, w router.ResponseWriter, r *router.Request) {
    w.SetStatus(500)
    w.WriteString(e.Error())
}

func myErrorHandler(e error, w router.ResponseWriter, r *router.Request) {
    if err, ok := e.(MyError); ok {
        w.SetStatus(err.HttpStatus)
        w.WriteAsJson(ResponseObject{Status: Status{Code: err.StatusCode , Description: err.Description}, Data: nil})
    } else {
        w.SetStatus(500)
        w.WriteAsJson(ResponseObject{Status: Status{Code: 500 , Description: e.Error()}, Data: nil})
    }
}

// Controller

func helloWorld(w router.ResponseWriter, r *router.Request) error {
    return w.Write([]byte("Hello World"))
}

func helloJson(w router.ResponseWriter, r *router.Request) error {
    res := ResponseObject{}

    res.Status.Code = 200
    res.Status.Description = "success"
    res.Data = "Hello Json"

    return w.WriteAsJson(res)
}

func helloMyError(w router.ResponseWriter, r *router.Request) error {
    return MyError{400, 9000, "Hello My Error"}
}

func helloError(w router.ResponseWriter, r *router.Request) error {
    return errors.New("Hello Error")
}
