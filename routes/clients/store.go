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

	tx, err := c.store.Begin(ctx)

	if err != nil {
		log.Println("error creating pg transaction")
		return err
	}

	defer func(){
		tx.Rollback(ctx)
	}()

	insertQuery := `
		INSERT INTO users (name, email, phone_number, pwd_hash, profile_pic)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`

	var user_id string
	err = tx.QueryRow(ctx, insertQuery, client.Name, client.Email, client.PhoneNumber,
		client.Password, client.ProfilePic).Scan(&user_id)

	if err != nil {
		return err
	}

	insertWalletQuery := `INSERT INTO wallets (user_id) VALUES ($1)`

	conn, err := tx.Exec(ctx, insertWalletQuery, user_id)

	if err != nil {
		log.Println(conn.String())
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("error while committing transaction")
		return err
	}

	return nil
}