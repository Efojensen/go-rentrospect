package assets

import (
	"encoding/json"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (a *AssetHandler) UploadAsset (w http.ResponseWriter, r *http.Request) {
	var newAsset types.Asset
	if err := json.NewDecoder(r.Body).Decode(&newAsset); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
	}

	a.addAssetQuery(newAsset)
}