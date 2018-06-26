package controller

import (
	"encoding/json"
	"net/http"

	entity "github.com/blanvam/rasp-garden/entities"
	"github.com/blanvam/rasp-garden/utils"
)

//ContextKeyResource is a key used for indexing a resource in a context
var ContextKeyResource = utils.ContextKey("resource")

func getResourceRequest(r *http.Request) (entity.ResourceRequest, error) {
	var resourceRequest entity.ResourceRequest
	err := json.NewDecoder(r.Body).Decode(&resourceRequest)
	return resourceRequest, err
}
