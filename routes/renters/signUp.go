package renters

import (
	"encoding/json"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (rr *RenterHandler) RenterSignUp(w http.ResponseWriter, r *http.Request) {
	var renterBody types.Renter

	if err := json.NewDecoder(r.Body).Decode(&renterBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	if err := rr.addRenterQuery(renterBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(w, http.StatusCreated, renterBody)
}