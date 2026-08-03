package clients

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (c *ClientHandler) getCategorizedAssets(w http.ResponseWriter, r *http.Request) {
	categoryId := r.Header.Get("categoryId")

	if categoryId == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, errors.New("missing category id"))
		return
	}

	returnedAssets, err := c.getCategorizedAssetsQuery(categoryId)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	for i, asset := range returnedAssets {
		url, err := c.cache.GetURL(asset.PrimaryImage)

		if err != nil {
			log.Printf("failed to generate url for index %d", i)
			url = ""
		}

		returnedAssets[i].PrimaryImage = url
	}

	if err = json.NewEncoder(w).Encode(returnedAssets); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *ClientHandler) getCategorizedAssetsQuery(categoryId string) ([]types.PureAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	categorizedAssets := make([]types.PureAsset, 0, 10)

	getAssetsByCategoryQuery := `
		SELECT a.vendor_id AS vendor, c.category_name AS category,a.asset_name AS name,
			i.object_name AS primaryImage, a.rate, a.pricing_unit, a.location,
			a.condition
			FROM assets AS a
		INNER JOIN categories AS c ON a.category_id = c.category_id
		INNER JOIN asset_images AS i ON a.asset_id = i.asset_id
		WHERE a.category_id = $1 AND availability_status = 'available'
		AND i.is_primary = TRUE
		LIMIT 10
	`

	rows, err := c.store.Query(ctx, getAssetsByCategoryQuery, categoryId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var queriedAsset types.PureAsset

		err := rows.Scan(
			&queriedAsset.Vendor,
			&queriedAsset.Category,
			&queriedAsset.Name,
			&queriedAsset.PrimaryImage,
			&queriedAsset.Rate,
			&queriedAsset.PricingUnit,
			&queriedAsset.Location,
			&queriedAsset.Condition,
		)

		if err != nil {
			log.Printf("Error scanning asset row: %v", err)
			continue
		}
		categorizedAssets = append(categorizedAssets, queriedAsset)
	}

	return categorizedAssets, nil
}
