package model

import "testing"

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet()

	if len(wallet.PrivKey) == 0 || len(wallet.PubKey) == 0 || err != nil {
		t.Error("Wallet not created:", wallet, err)
	}
}
