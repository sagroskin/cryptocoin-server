package service

import (
	"cryptocoin-server/model"
	"testing"
	"time"
)

func TestVerifyTranscation(t *testing.T) {
	pw, _ := model.NewWallet()
	nw, _ := model.NewWallet()

	pt := new(model.Transaction)

	pt.Value = 100
	pt.Timestamp = time.Now()
	pt.PubKey = pw.PubKey
	pt.ToAddress = nw.PubKey
	pt.Signature, _ = CalculateSignature(pt, pw.PrivKey)

	nt := new(model.Transaction)

	nt.PrevSignature = pt.Signature
	nt.Value = 100
	nt.Timestamp = time.Now()
	nt.PubKey = nw.PubKey
	nt.ToAddress = nw.PubKey // Send to send
	nt.Signature, _ = CalculateSignature(nt, nw.PrivKey)

	result, err := VerifyTransaction(nt, pt)

	if err != nil || result != true {
		t.Error("VerifyTransaction failed:", nt, pt, err)
	}
}

func TestCalculateSignature(t *testing.T) {
	wallet, err := model.NewWallet()
	transacation := new(model.Transaction)

	transacation.Value = 100
	transacation.Timestamp = time.Now()
	transacation.PubKey = wallet.PubKey
	transacation.Signature, _ = CalculateSignature(transacation, wallet.PrivKey)

	if err != nil {
		t.Error("CalculateSignature failed:", transacation, wallet)
	}
}

func TestVerifySignature(t *testing.T) {
	wallet, err := model.NewWallet()
	transacation := new(model.Transaction)

	transacation.Value = 100
	transacation.Timestamp = time.Now()
	transacation.PubKey = wallet.PubKey
	transacation.Signature, _ = CalculateSignature(transacation, wallet.PrivKey)

	result, err := VerifySignature(transacation)

	if err != nil || result == false {
		t.Error("TestVerifySignature failed:", "transaction=", transacation, "wallet=", wallet)
	}
}
