package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
)

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error, msg ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	code := types.ErrorCodeEnum(statusCode)

	apiError := types.CustomError{
		Code:    code.ToString(),
	}

	if len(msg) > 0 {
		apiError.Message = &msg[0]
	}

	err = json.NewEncoder(w).Encode(apiError)

	if err != nil {
		log.Println("failed to encode error response: ", err)
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
