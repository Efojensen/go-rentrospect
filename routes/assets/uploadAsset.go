package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (a *AssetHandler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	assetJson := r.FormValue("assetDetails")
	if assetJson == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, errors.New("missing event data"))
		return
	}

	var assetInfo types.Asset
	if err := json.Unmarshal([]byte(assetJson), &assetInfo); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, errors.New("invalid event json"))
		return
	}

	assetImages := make([]types.AssetImage, 0, 4)

	for fieldName, fileHeaders := range r.MultipartForm.File {
		for i, header := range fileHeaders {
			fmt.Printf("Processing field: %s, file index: %d, name: %s\n", fieldName, i, header.Filename)

			// 4. Open the file
			file, err := header.Open()
			if err != nil {
				http.Error(w, "Failed to read file", http.StatusInternalServerError)
				return
			}

			defer file.Close()

			// 5. Validate that it is actually an image (Sniff content type)
			fileBytes, err := io.ReadAll(file)

			contentType := http.DetectContentType(fileBytes)
			if !utils.IsImage(contentType) {
				break
			}
			var assetImage types.AssetImage = types.AssetImage{
				FileBytes:   fileBytes,
				ContentType: contentType,
			}
			assetImages = append(assetImages, assetImage)
		}
	}

	if err := a.addAssetQuery(assetInfo, assetImages); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
}
