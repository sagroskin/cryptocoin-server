package main

import (
	"fmt"
	"log"
	"net/http"
	"cryptocoin-server/config"
	"cryptocoin-server/controller"

	"github.com/gorilla/mux"
)

func main() {
	config := config.InitConfig()
	router := mux.NewRouter()

	controller.InitTransactionController(router)
	controller.InitWalletController(router)

	fmt.Println("Started server http://localhost" + config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}
