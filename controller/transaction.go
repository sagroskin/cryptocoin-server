package controller

import (
	"cryptocoin-server/model"
	"cryptocoin-server/repository"
	"cryptocoin-server/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// InitTransactionController initializes the controller.
func InitTransactionController(router *mux.Router) {
	router.HandleFunc("/transactions", GetTransactions).Methods("GET")
	router.HandleFunc("/transactions", CreateTransactions).Methods("POST")
	router.HandleFunc("/transactions/genesis", CreateGenesisTransaction).Methods("POST")
	router.HandleFunc("/transactions/transfer", TransferFromGenesisAccount).Methods("POST")
	router.HandleFunc("/transactions/{id}", GetTransaction).Methods("GET")
}

// GetTransactions returns all latest transactions (limit 25).
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := repository.GetTransactions()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

// CreateTransactions submit new transactions to create. All transactions must be valid.
func CreateTransactions(w http.ResponseWriter, r *http.Request) {
	// Get Transaction from body
	var transactions []model.Transaction
	err := json.NewDecoder(r.Body).Decode(&transactions)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Add Transaction
	service.AddTransactions(transactions)

	// Return confirmed Transaction
	json.NewEncoder(w).Encode(transactions)
}

// GetTransaction returns transaction by ID (Signature).
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	signature := params["id"]

	transaction, err := repository.GetTransaction(signature, true)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

// GetTransactionHistory returns the current transaction and the history (all previous transactions).
func GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	// signature := params["id"]

	transactions, err := repository.GetTransactions() //TODO: Change to get history for specific transaction

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

// CreateGenesisTransaction creates a new genesis transaction. For testing use only.
func CreateGenesisTransaction(w http.ResponseWriter, r *http.Request) {
	t, err := service.CreateGenesisTransaction()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(t)
}

// TransferFromGenesisAccount transfer a specified amount to a wallet from the Genesis wallet. For testing use only.
func TransferFromGenesisAccount(w http.ResponseWriter, r *http.Request) {
	signature := r.FormValue("signature")
	sendTo := r.FormValue("sendTo")
	amount, _ := strconv.ParseInt(r.FormValue("amount"), 10, 64)

	transactions, err := service.TransferFromGenesisAccount(signature, sendTo, amount)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
