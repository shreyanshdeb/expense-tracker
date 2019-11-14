package main

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JSONExpense struct {
	ID       string  `json:"ID"`
	Title    string  `json:"Title"`
	Category string  `json:"Category"`
	Amount   float64 `json:"Amount"`
	DateTime int64   `json:"Datetime"`
	UserID   string  `json:"Userid"`
}

type Expense struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"datetime"`
	UserID   string    `json:"userid"`
}

type ExpenseBreakdown struct {
	Savings []*Expense
	Needs   []*Expense
	Wants   []*Expense
}

type User struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Email    string `json: "Email"`
	Phone    string `json: "Phone"`
	Password string `json: "Password"`
}

type Claims struct {
	UserID string
	Email  string
	jwt.StandardClaims
}

var Expensedb = []*Expense{}

var Userdb = []*User{}

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
