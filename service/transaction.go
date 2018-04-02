package service

import (
	"cryptocoin-server/config"
	"cryptocoin-server/model"
	"cryptocoin-server/repository"
	"cryptocoin-server/util/ecdsa"
	"errors"
	"time"
)

// AddTransactions adds new transactions to the ledger. All transactions must be valid.
func AddTransactions(transactions []model.Transaction) error {
	prevTransactions := make(map[string]*model.Transaction)

	// Get previous transaction for each new transaction
	for _, t := range transactions {
		if _, contains := prevTransactions[t.PrevSignature]; !contains {
			pt, err := repository.GetTransaction(t.PrevSignature, false) // TODO: Get all previous transactions in one call

			if err != nil {
				return err
			}

			if pt != nil {
				prevTransactions[pt.Signature] = pt
			}
		}
	}

	timestamp := transactions[0].Timestamp

	// Verify all the transactions and calculate new balance
	for _, t := range transactions {
		verified, err := VerifyTransaction(&t, prevTransactions[t.PrevSignature])

		if err != nil {
			return err
		}

		if !verified {
			return errors.New("At least one of the transactions was invalid")
		}

		if !t.Timestamp.Equal(timestamp) {
			return errors.New("Timestamp must be same for all transactions")
		}

		prevTransactions[t.PrevSignature].Value -= t.Value
	}

	// Verify that all previous transactions will now have 0 balance
	for _, pt := range prevTransactions {
		if pt.Value != 0 {
			return errors.New("Total value of all transactions (in and out) must be equal")
		}
	}

	// Save transactions
	for _, t := range transactions {
		err := repository.SaveTransaction(&t) // TODO: Save all transactions in one call at the end

		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyTransaction verifies all properties for the transaction.
func VerifyTransaction(t *model.Transaction, pt *model.Transaction) (bool, error) {
	tv, _ := VerifySignature(t)
	ptv, _ := VerifySignature(t)

	// Verify prevTrasaction was found
	if pt == nil {
		return false, errors.New("Previous transaction not found or is invalid")
	}

	// Verify Signature of both Transactions
	if !tv || !ptv {
		return false, errors.New("Signatures must be valid for both transactions")
	}

	// Verify that new transaction matches prevTransaction signature
	if t.PrevSignature != pt.Signature {
		return false, errors.New("Signature must match previous transaction")
	}

	// Verify that the owner matches
	if pt.ToAddress != t.PubKey {
		return false, errors.New("Owner must match previous transaction")
	}

	// Verify Timestamp of new transaction > prevTransaction
	if pt.Timestamp.After(t.Timestamp) {
		return false, errors.New("Timestamp must be after previous transaction")
	}

	// Verify Timestamp not in the future
	if t.Timestamp.After(time.Now()) {
		return false, errors.New("Timestamp must be less than current time")
	}

	// Verify SendTo is a valid Public Key
	_, err := ecdsa.ParsePubKey(t.ToAddress)
	if err != nil {
		return false, errors.New("SendTo must be a valid Public Key")
	}

	// Verify prevTrasaction > new transaction and subtract from it
	if t.Value > pt.Value {
		return false, errors.New("Previous transaction value must be equal to the combined value of new transactions")
	}

	return true, nil
}

// VerifySignature verifies the signature of a transaction using the transaction hash and public key.
func VerifySignature(transaction *model.Transaction) (bool, error) {
	key, err := ecdsa.ParsePubKey(transaction.PubKey)
	if err != nil {
		return false, err
	}

	hash, err := transaction.Hash()
	if err != nil {
		return false, err
	}

	result, err := ecdsa.Verify(key, hash[:], transaction.Signature)

	if err != nil {
		return false, err
	}

	return result, nil
}

// CalculateSignature calculates the signature for a transaction using the hash of the transaction and private key.
func CalculateSignature(transaction *model.Transaction, privKey string) (string, error) {
	hash, err := transaction.Hash()

	if err != nil {
		return "", err
	}

	key, err := ecdsa.ParsePrivKey(privKey)

	if err != nil {
		return "", err
	}

	signature, err := ecdsa.Sign(key, hash[:])

	if err != nil {
		return "", err
	}

	return signature, nil
}

// CreateGenesisTransaction creates a new genesis transaction. For testing use only.
func CreateGenesisTransaction() (*model.Transaction, error) {
	config := config.InitConfig()
	t := model.NewTransaction()

	t.Value = 1000000
	t.PrevSignature = "GENESIS"
	t.PubKey = config.GenesisPubKey
	t.ToAddress = config.GenesisPubKey
	t.Timestamp = time.Now()

	signature, err := CalculateSignature(t, config.GenesisPrivKey)
	if err != nil {
		return nil, err
	}

	t.Signature = signature

	err = repository.SaveTransaction(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// TransferFromGenesisAccount transfer a specified amount to a wallet from the Genesis wallet. For testing use only.
func TransferFromGenesisAccount(signature string, sendTo string, amount int64) (*[]model.Transaction, error) {
	config := config.InitConfig()

	pt, err := repository.GetTransaction(signature, true)
	if err != nil {
		return nil, err
	}

	if pt == nil {
		return nil, errors.New("Transaction does not exist")
	}

	timestamp := time.Now()

	transactions := make([]model.Transaction, 2)
	transactions[0].Value = amount
	transactions[0].PubKey = config.GenesisPubKey
	transactions[0].Timestamp = timestamp
	transactions[0].PrevSignature = pt.Signature
	transactions[0].ToAddress = sendTo
	signature1, err := CalculateSignature(&transactions[0], config.GenesisPrivKey)
	if err != nil {
		return nil, err
	}
	transactions[0].Signature = signature1

	transactions[1].Value = pt.Value - amount
	transactions[1].PubKey = config.GenesisPubKey
	transactions[1].Timestamp = timestamp
	transactions[1].PrevSignature = pt.Signature
	transactions[1].ToAddress = config.GenesisPubKey
	signature2, err := CalculateSignature(&transactions[1], config.GenesisPrivKey)
	if err != nil {
		return nil, err
	}
	transactions[1].Signature = signature2

	err = AddTransactions(transactions)

	return &transactions, err
}
