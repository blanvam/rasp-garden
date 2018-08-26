package api

import "net/http"

// Controller interface definition
type Controller interface {
	ListHandler(res http.ResponseWriter, req *http.Request)
	CreateHandler(res http.ResponseWriter, req *http.Request)
	DetailHandler(res http.ResponseWriter, req *http.Request)
	DeleteHandler(res http.ResponseWriter, req *http.Request)
}
