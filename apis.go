package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func getExpensesHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, _ := json.Marshal(getExpenses())
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	expense, err := getExpenseByID(id)
	if err != nil {
		jsonData, _ := json.Marshal(`{
			"sucess":false,
			"message": "Expense not found"
		}`)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
	jsonData, _ := json.Marshal(expense)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func addExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var je JSONExpense
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()

	id, _ := uuid.NewV4()
	expense.ID = id.String()
	addExpense(*expense)

	jsonData, _ := json.Marshal(expense)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	deleted := deleteExpenseByID(id)
	if !deleted {
		jsonData, _ := json.Marshal(`{"success":false}`)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}
	jsonData, _ := json.Marshal(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func updateExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	var je JSONExpense

	params := mux.Vars(r)
	id := params["id"]
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()

	updated := updateExpenseByID(id, *expense)

	if !updated {
		jsonData, _ := json.Marshal(`{"sucess" : false}`)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	jsonData, _ := json.Marshal(`{"sucess" : true}`)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}
