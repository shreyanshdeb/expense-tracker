package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func initDb() {
	Expensedb = append(Expensedb, &Expense{
		1,
		"Rent",
		"Needs",
		4500,
		"2019-10-01",
		"",
	})
	Expensedb = append(Expensedb, &Expense{
		2,
		"House Help",
		"Needs",
		5000,
		"2019-10-01",
		"",
	})
}

func main() {
	initDb()

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/expenses", getExpensesHandler).Methods("GET")
	r.HandleFunc("/api/v1/expense/{id}", getExpenseByIDHandler).Methods("GET")
	r.HandleFunc("/api/v1/expense", addExpenseHandler).Methods("POST")
	r.HandleFunc("/api/v1/expense/{id}", deleteExpenseByIDHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/expense/{id}", updateExpenseByIDHandler).Methods("PUT")

	fmt.Println("Listening on 1234")
	log.Fatal(http.ListenAndServe(":1234", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
