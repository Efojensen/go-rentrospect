package renters

import (
	"context"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (r *RenterHandler) addRenterQuery(renter types.Renter) error {
	ctx, timeout := context.WithTimeout(context.Background(), 5 * time.Second)
	defer timeout()

	insertQuery := `
		INSERT INTO users (name, email, phone_number, password, profile_pic)
		VALUES ($1, $2, $3, $4, $5)
	`

	commandTag, err := r.store.Exec(ctx, insertQuery, renter.Name, renter.Email, renter.PhoneNumber,
		renter.Password, renter.ProfilePic)

	if err != nil {
		return err
	}

	log.Println(commandTag.RowsAffected(), "rows affected")
	return nil
}