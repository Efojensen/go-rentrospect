package assets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
	"github.com/EfoJensen/go-rentrospect/utils"
)

func (a *AssetHandler) addAssetQuery(asset types.Asset, assetImages []types.AssetImage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 14*time.Second)
	defer cancel()

	objectNames := make([]string, 0, 4)
	if asset.PrimaryImage < 0 || asset.PrimaryImage >= len(assetImages) {
		return errors.New("invalid primary image index")
	}

	for _, file := range assetImages {
		ext := utils.GetExtensionFromContentType(file.ContentType)

		objectName := fmt.Sprintf(
			"assets/%d-%s%s",
			time.Now().Unix(),
			utils.GenerateRandomString(8),
			ext,
		)

		_, err := a.objectStore.UploadFromBytes(objectName, file.FileBytes)
		if err != nil {
			return err
		}
		objectNames = append(objectNames, objectName)
	}

	tx, err := a.store.Begin(ctx)
	if err != nil {
		log.Println("failed to begin transaction: %w", err)
		return err
	}

	// Database rollback in case of failed transaction. Perfectly fine to have defer if
	// Commit was called before the defer
	defer func() {
		tx.Rollback(ctx)
	}()

	defer func() {
		if err != nil {
			for _, obj := range objectNames {
				if delErr := a.objectStore.DeleteObject(obj); delErr != nil {
					log.Printf("failed to delete orphaned object %q: %v", obj, delErr)
				}
			}
		}
	}()

	var createdAssetID string
	insertQuery := `
		INSERT INTO assets
		(vendor_id, category_id, asset_name, availability_status, description, rate,
			pricing_unit, location, condition)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING asset_id;
	`
	err = tx.QueryRow(ctx, insertQuery, asset.Vendor, asset.Category, asset.Name,
		asset.Availability, asset.Description, asset.Rate, asset.PricingUnit,
		asset.Location, asset.Condition,
	).Scan(&createdAssetID)

	if err != nil {
		return err
	}

	for idx, objName := range objectNames {
		isPrimary := false
		if idx == asset.PrimaryImage {
			isPrimary = true
		}

		imagesInsertQuery := `
			INSERT INTO asset_images (asset_id, object_name, is_primary)
			VALUES ($1, $2, $3)
		`
		if _, err := tx.Exec(ctx, imagesInsertQuery, createdAssetID, objName, isPrimary); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
