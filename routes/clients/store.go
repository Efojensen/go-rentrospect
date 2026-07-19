package clients

import (
	"context"
	"log"
	"time"

	"github.com/EfoJensen/go-rentrospect/types"
)

func (c *ClientHandler) addClientQuery(client types.Client) error {
	ctx, timeout := context.WithTimeout(context.Background(), 5 * time.Second)
	defer timeout()

	insertQuery := `
		INSERT INTO users (name, email, phone_number, password, profile_pic)
		VALUES ($1, $2, $3, $4, $5)
	`

	commandTag, err := c.store.Exec(ctx, insertQuery, client.Name, client.Email, client.PhoneNumber,
		client.Password, client.ProfilePic)

	if err != nil {
		return err
	}

	log.Println(commandTag.RowsAffected(), "rows affected in \"users\" table")
	return nil
}