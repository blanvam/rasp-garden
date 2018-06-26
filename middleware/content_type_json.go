package middleware

import (
	"log"
	"net/http"

	"github.com/blanvam/rasp-garden/utils"
)

// RequireJSONMiddleware is a struct that has a ServeHTTP method
type RequireJSONMiddleware struct {
}

// NewRequireJSONMiddleware returns a RequireJSONMiddleware
func NewRequireJSONMiddleware() *RequireJSONMiddleware {
	return &RequireJSONMiddleware{}
}

/*
RequireJSONMiddleware's handler, which asserts that POST and PUT methods include content-type header
and is set to application/json
*/
func (l *RequireJSONMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	methodNeedsJSON := func(method string) bool {
		return method == "POST" || method == "PUT"
	}
	if methodNeedsJSON(r.Method) && r.Header.Get("content-type") != "application/json" {
		log.Println("content-type bust be 'application/json'")
		utils.ResponseError(w, "Expected content-type to be application/json", http.StatusBadRequest)
	} else {
		next(w, r)
	}
}
