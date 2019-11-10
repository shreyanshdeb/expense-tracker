package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func initDb() {
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	Expensedb = append(Expensedb, &Expense{
		id1.String(),
		"Rent",
		"Needs",
		4500,
		time.Unix(1480979203, 0),
		"",
	})
	Expensedb = append(Expensedb, &Expense{
		id2.String(),
		"House Help",
		"Needs",
		5000,
		time.Unix(1480979203, 0),
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
