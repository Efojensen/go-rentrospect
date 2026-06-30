package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
)

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	apiError := types.Errors{
		Error:   err.Error(),
	}

	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	log.Println(err)
}

func WriteResponse(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(msg); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
}
