package utils

import (
	"encoding/json"
	"net/http"
)

type listSerializer struct {
	Results interface{} `json:"results"`
}

// ContextKey is a string used for key indexing at for context
type ContextKey string

//ContextKeyResource is a key used for indexing a resource in a context
var ContextKeyResource = ContextKey("resource")

//RemoveForbiddenFields removes id created_at and modified at from JSONMap
func getJSONEncoder(w http.ResponseWriter) *json.Encoder {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder
}

// ResponseError returns writes to w a mesage and sets the status code to code
func ResponseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	getJSONEncoder(w).Encode(map[string]string{"error": message})
}

// ResponseJSON serializes a object and sends the result to w
func ResponseJSON(w http.ResponseWriter, object interface{}, many bool) { //serializer serializerFn, object interface{}) {
	var err error

	w.Header().Set("Content-Type", "application/json")
	encoder := getJSONEncoder(w)
	if many {
		err = encoder.Encode(listSerializer{Results: object})
	} else {
		err = encoder.Encode(object)
	}
	if err != nil {
		ResponseError(w, err.Error(), http.StatusBadRequest)
	}
}

// ResponseCreated sets header to 201 Created
func ResponseCreated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

// ResponseNoContent sets header to 204 NoContent
func ResponseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
