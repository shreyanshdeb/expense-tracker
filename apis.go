package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getExpensesHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, _ := json.Marshal(getExpenses())
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	intID, _ := strconv.ParseInt(id, 10, 32)
	expense, err := getExpenseByID(intID)
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
	var expense Expense
	_ = json.NewDecoder(r.Body).Decode(&expense)

	expense.ID = int64(len(Expensedb) + 1)
	addExpense(expense)

	jsonData, _ := json.Marshal(expense)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 32)
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
	var expense Expense

	params := mux.Vars(r)
	id := params["id"]
	idInt, _ := strconv.ParseInt(id, 10, 64)
	_ = json.NewDecoder(r.Body).Decode(&expense)

	updated := updateExpenseByID(idInt, expense)

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
