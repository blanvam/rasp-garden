package api

import "net/http"

type Middleware interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}
