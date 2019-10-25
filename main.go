package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Expense struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	DateTime    string  `json:"datetime"`
	UserID      string  `json:"userid"`
}

var Expensedb = []*Expense{}

var ErrExpenseNotFound = errors.New("Expense not found.")

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
	r.HandleFunc("/api/v1/expense/{id}", deleteExpenseByIDHandler).Methods("Delete")

	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", r)
}

func getExpenses() []Expense {
	AllExpenses := []Expense{}
	for _, v := range Expensedb {
		AllExpenses = append(AllExpenses, *v)
	}
	return AllExpenses
}

func addExpense(expense Expense) {
	Expensedb = append(Expensedb, &expense)
}

func getExpenseByID(id int64) (*Expense, error) {
	for _, v := range Expensedb {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, ErrExpenseNotFound
}

func deleteExpenseByID(id int64) bool {
	for i, v := range Expensedb {
		if v.ID == id {
			Expensedb = append(Expensedb[:i], Expensedb[i+1:]...)
			return true
		}
	}
	return false
}

func updateExpenseByID(id int64, expense Expense) bool {
	for _, v := range Expensedb {
		if v.ID == id {
			v.Title = expense.Title
			v.Description = expense.Title
			v.Amount = expense.Amount
			v.DateTime = expense.DateTime
			v.UserID = expense.UserID
			return true
		}
	}
	return false
}

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
	var expense *Expense
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), expense); err != nil {
			fmt.Fprintln(w, err)
			jsonData, _ := json.Marshal(`{"success":false}`)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}
	}
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

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
