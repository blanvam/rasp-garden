package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/blanvam/rasp-garden/resource"
	"github.com/blanvam/rasp-garden/resource/controller"
	"github.com/blanvam/rasp-garden/utils"
	"github.com/gorilla/mux"
)

//RequireResourceMiddleware is a middleware that ensures a url's id parameter is a valid ID related to a Resource document
type RequireResourceMiddleware struct {
	usecase resource.Usecase
}

// NewRequireResourceMiddleware returns a RequireResourceMiddleware
func NewRequireResourceMiddleware(u resource.Usecase) *RequireResourceMiddleware {
	return &RequireResourceMiddleware{
		usecase: u,
	}
}

func getResourceID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	return id, nil
}

/*
RequireResourceMiddleware's handler, which asserts that url's id parameter is a valid ID and is related to a Resource
document in the database
*/
func (l *RequireResourceMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	log.Println("Middleware resource serveHTTP")
	resourceID, _ := getResourceID(req)
	log.Println("Middleware got id: " + strconv.Itoa(resourceID))
	ctx := req.Context()
	resource, err := l.usecase.GetByID(ctx, resourceID)
	if err != nil {
		utils.ResponseError(res, err.Error(), http.StatusNotFound)
		return
	}

	req = req.WithContext(context.WithValue(req.Context(), controller.ContextKeyResource, resource))

	next(res, req)
}
