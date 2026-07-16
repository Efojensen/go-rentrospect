package vendors

import (
	"encoding/json"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (v *VendorHandler) VendorSignUp(w http.ResponseWriter, r *http.Request) {
	var vendorBody types.Vendor

	if err := json.NewDecoder(r.Body).Decode(&vendorBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := v.storeVendorQuery(vendorBody); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
}