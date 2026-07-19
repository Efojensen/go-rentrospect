package clients

import (
	"encoding/json"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (c *ClientHandler) ClientSignUp(w http.ResponseWriter, r *http.Request) {
	var clientBody types.Client

	if err := json.NewDecoder(r.Body).Decode(&clientBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	if err := c.addClientQuery(clientBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(w, http.StatusCreated, clientBody)
}