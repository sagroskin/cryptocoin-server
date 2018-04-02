package controller

import (
	"cryptocoin-server/model"
	"cryptocoin-server/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// InitWalletController initializes the controller.
func InitWalletController(router *mux.Router) {
	router.HandleFunc("/wallet", GetWallet).Methods("GET")
	router.HandleFunc("/wallet/create", CreateWallet).Methods("GET")
}

// GetWallet returns the wallet by public key.
func GetWallet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pubKey := params["id"]

	wallet, err := service.GetWallet(pubKey)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = json.NewEncoder(w).Encode(wallet)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

// CreateWallet creates a new wallet with a private and public key.
func CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet, err := model.NewWallet()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = json.NewEncoder(w).Encode(wallet)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}
