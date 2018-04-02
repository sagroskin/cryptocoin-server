package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Transaction struct
type Transaction struct {
	Timestamp     time.Time `json:"timestamp"`
	ToAddress     string    `json:"toAddress"`
	Value         int64     `json:"value"`
	PubKey        string    `json:"pubKey"`
	PrevSignature string    `json:"prevSignature"`
	Signature     string    `json:"signature"`
}

// NewTransaction ?
func NewTransaction() *Transaction {
	t := new(Transaction)
	return t
}

// Serialize ?
func (transaction *Transaction) Serialize() ([]byte, error) {
	var data bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&data)

	err := enc.Encode(transaction.Timestamp)

	if err != nil {
		return nil, err
	}

	err = enc.Encode(transaction.ToAddress)

	if err != nil {
		return nil, err
	}

	err = enc.Encode(transaction.Value)

	if err != nil {
		return nil, err
	}

	err = enc.Encode(transaction.PubKey)

	if err != nil {
		return nil, err
	}

	err = enc.Encode(transaction.PrevSignature)

	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

// Hash ?
func (transaction *Transaction) Hash() ([]byte, error) {
	data, err := transaction.Serialize()

	if err != nil {
		return nil, err
	}

	hashBytes := sha256.Sum256(data)

	return hashBytes[:], nil
}
