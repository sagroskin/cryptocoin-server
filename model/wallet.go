package model

import (
	"cryptocoin-server/util/ecdsa"
)

// Wallet struct containing private key, public key, and balance
type Wallet struct {
	PrivKey   string
	PubKey    string
	Balanance int64
}

// NewWallet creates a new wallet with a private and public key.
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateNewKey()
	if err != nil {
		return nil, err
	}

	w := new(Wallet)

	w.PrivKey = ecdsa.ExportPrivKey(privateKey)
	w.PubKey = ecdsa.ExportPubKey(&privateKey.PublicKey)
	w.Balanance = 0

	return w, nil
}
