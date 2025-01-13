package main

import (
	"log"
	"net/http"

	"github.com/alvinatthariq/billing-engine/internal/api"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	loanHandler := api.NewLoanHandler()

	// Register routes
	router.HandleFunc("/loans", loanHandler.CreateLoan).Methods("POST")
	router.HandleFunc("/loans/{id}/outstanding", loanHandler.GetOutstanding).Methods("GET")
	router.HandleFunc("/loans/{id}/delinquent", loanHandler.IsDelinquent).Methods("GET")
	router.HandleFunc("/loans/{id}/payment", loanHandler.MakePayment).Methods("POST")

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}
