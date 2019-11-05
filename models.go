package main

type Expense struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Category 	string  `json:"category"`
	Amount      float64 `json:"amount"`
	DateTime    string  `json:"datetime"`
	UserID      string  `json:"userid"`
}

var Expensedb = []*Expense{}
