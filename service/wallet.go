package service

import (
	"cryptocoin-server/model"
)

// GetWallet returns wallet
func GetWallet(pubKey string) (*model.Wallet, error) {
	wallet := new(model.Wallet) //TODO: Get wallet summary
	return wallet, nil
}
