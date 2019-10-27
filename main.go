package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func initDb() {
	Expensedb = append(Expensedb, &Expense{
		1,
		"Rent",
		"Rent for the month October",
		4500,
		"2019-10-01",
		"",
	})
	Expensedb = append(Expensedb, &Expense{
		2,
		"House Help",
		"House Help Salary for the month October",
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

	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", r)
}
