package controller

import (
	"context"
	"net/http"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/resource"
	"github.com/blanvam/rasp-garden/utils"
)

// ResponseError notify the cause of a error
type ResponseError struct {
	Message string `json:"message"`
}

// HTTPResourceHandler handler
type HTTPResourceHandler struct {
	usecase resource.Usecase
}

// NewResourceHTTPpHandler aaa
func NewResourceHTTPpHandler(u resource.Usecase) *HTTPResourceHandler {
	return &HTTPResourceHandler{
		usecase: u,
	}
}

// ListResourcesHandler handles GET requests for Resource listing
func (r *HTTPResourceHandler) ListResourcesHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resources, err := r.usecase.All(ctx)
	if err != nil {
		utils.ResponseError(res, err.Error(), http.StatusBadRequest)
	} else {
		utils.ResponseJSON(res, resources, true)
	}
}

// CreateResourceHandler handles POST requests for resource creation
func (r *HTTPResourceHandler) CreateResourceHandler(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resourceRequest, err := getResourceRequest(req)
	if err != nil {
		utils.ResponseError(res, err.Error(), http.StatusBadRequest)
		return
	}

	memResource, _ := r.usecase.Bind(ctx, &resourceRequest)
	stored, err := r.usecase.Store(ctx, memResource)

	if stored == false {
		utils.ResponseError(res, err.Error(), http.StatusBadRequest)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
}

// ResourceDetailHandler handles GET requests for resource detail
func (r *HTTPResourceHandler) ResourceDetailHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	resource := ctx.Value(ContextKeyResource).(*entity.Resource)
	utils.ResponseJSON(res, resource, false)
}

// ResourceDeleteHandler handles DELETE requests for resource deletion
func (r *HTTPResourceHandler) ResourceDeleteHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	resource := ctx.Value(ContextKeyResource).(*entity.Resource)

	deleted, err := r.usecase.Delete(ctx, resource.Pin)
	if err != nil {
		utils.ResponseError(res, err.Error(), http.StatusBadRequest)
		return
	}
	if deleted == false {
		utils.ResponseError(res, "Resource can not be deleted", http.StatusBadRequest)
	}
	res.WriteHeader(http.StatusNoContent)
}
