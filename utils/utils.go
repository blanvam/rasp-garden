package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ContextKey is a string used for key indexing at for context
type ContextKey string

type listSerializer struct {
	Results interface{} `json:"results"`
}

func getJSONEncoder(w http.ResponseWriter) *json.Encoder {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder
}

func badRequestHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ResponseError(res, "Expected content-type to be application/json", http.StatusBadRequest)
	})
}

// RequireJSON is a MatcherFunc for gorilla mux, which specifies that a method is accesed with json
func RequireJSON(req *http.Request, rm *mux.RouteMatch) bool {
	if req.Method == "POST" && req.Header.Get("content-type") != "application/json" {
		rm.Handler = badRequestHandler()
	}
	return true
}

// ResponseError returns writes to w a mesage and sets the status code to code
func ResponseError(res http.ResponseWriter, message string, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	getJSONEncoder(res).Encode(map[string]string{"error": message})
}

// ResponseJSON serializes a object and sends the result to w
func ResponseJSON(res http.ResponseWriter, object interface{}, many bool) { //serializer serializerFn, object interface{}) {
	var err error

	res.Header().Set("Content-Type", "application/json")
	encoder := getJSONEncoder(res)
	if many {
		err = encoder.Encode(listSerializer{Results: object})
	} else {
		err = encoder.Encode(object)
	}
	if err != nil {
		ResponseError(res, err.Error(), http.StatusBadRequest)
	}
}
