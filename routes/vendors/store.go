package vendors

import (
	"context"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (v *VendorHandler) storeVendorQuery(vendor types.Vendor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)

	defer cancel()

	insertQuery := `
		INSERT INTO users (name, email, phone_number, password, profile_pic)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`

	tx, err := v.store.Begin(ctx)
	if err != nil {
		log.Println("failed to begin transaction: %w", err)
		return err
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	var user_id string

	err = tx.QueryRow(ctx, insertQuery, vendor.Name,
		vendor.Email, vendor.PhoneNumber, vendor.Password, vendor.ProfilePic,
	).Scan(&user_id)

	if err != nil {
		log.Println("error writing to \"users\" table: %w", err)
		return err
	}

	vendorInsertQuery := `
		INSERT INTO vendors (user_id, national_id)
		VALUES ($1, $2)
	`

	commandTag, err := tx.Exec(ctx, vendorInsertQuery, user_id, vendor.NatId)

	if err != nil {
		log.Println("error writing to \"vendors\" table: %w", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	log.Println(commandTag.RowsAffected(), "rows affected in \"vendors\" table")
	return nil
}
