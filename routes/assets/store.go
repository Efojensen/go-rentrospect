package assets

import (
	"context"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (a *AssetHandler) addAssetQuery(asset types.Asset) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertQuery := `
		INSERT INTO assets
		(vendor_id, category_id, asset_name, availability_status, description, rate,
			pricing_unit, location, condition)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	conn, err := a.store.Exec(ctx, insertQuery,
		asset.Vendor, asset.Category, asset.Name, asset.Availability, asset.Description,
		asset.Rate, asset.PricingUnit, asset.Location, asset.Condition,
	)

	if err != nil {
		return err
	}

	log.Printf("%d rows affected", conn.RowsAffected())

	return nil
}