package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/EfoJensen/go-rentrospect/types"
)

func GenerateIdempotencyKey(session types.PaymentSession) ([]byte, error) {
	data, err := json.Marshal(session)

	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}

	hash := sha256.Sum256(data)

	session.IdempotencyKey = hex.EncodeToString(hash[:])

	newData, err := json.Marshal(session)

	if err != nil {
		return nil, err
	}

	return newData, nil
}
