package main

import (
	"encoding/json"
	"time"
)

type JSONExpense struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	DateTime int64   `json:"datetime"`
	UserID   string  `json:"userid"`
}

type Expense struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"datetime"`
	UserID   string    `json:"userid"`
}

var Expensedb = []*Expense{}

func (je *JSONExpense) Expense() *Expense {
	return &Expense{
		je.ID,
		je.Title,
		je.Category,
		je.Amount,
		time.Unix(je.DateTime, 0),
		je.UserID,
	}
}

func NewJSONExpense(expense *Expense) *JSONExpense {
	return &JSONExpense{
		expense.ID,
		expense.Title,
		expense.Category,
		expense.Amount,
		expense.DateTime.Unix(),
		expense.UserID,
	}
}

func (e *Expense) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONExpense(e))
}

func (e *Expense) UnmarshalJSON(data []byte) ([]byte, error) {
	var je *JSONExpense
	if err := json.Unmarshal(data, je); err != nil {
		return nil, err
	}
	e = je.Expense()
	return nil, nil
}
