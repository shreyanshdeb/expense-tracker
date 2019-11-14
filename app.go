package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func getExpenses(userID string) []Expense {
	AllExpenses := []Expense{}
	for _, v := range Expensedb {
		if v.UserID == userID {
			AllExpenses = append(AllExpenses, *v)
		}
	}
	return AllExpenses
}

func addExpense(expense Expense) {
	Expensedb = append(Expensedb, &expense)
}

func getExpenseByID(id string, userID string) (*Expense, error) {
	for _, v := range Expensedb {
		if v.ID == id && v.UserID == userID {
			return v, nil
		}
	}
	return nil, ErrExpenseNotFound
}

func deleteExpenseByID(id string, userID string) bool {
	for i, v := range Expensedb {
		if v.ID == id && v.UserID == userID {
			Expensedb = append(Expensedb[:i], Expensedb[i+1:]...)
			return true
		}
	}
	return false
}

func updateExpenseByID(id string, expense Expense, userID string) bool {
	for _, v := range Expensedb {
		if v.ID == id && v.UserID == userID {
			v.Title = expense.Title
			v.Category = expense.Category
			v.Amount = expense.Amount
			v.DateTime = expense.DateTime
			return true
		}
	}
	return false
}

func listExpensesForMonth(date time.Time, userID string) []*Expense {
	var ExpenseForTheMonth = []*Expense{}
	for _, v := range Expensedb {
		if v.UserID == userID && (v.DateTime.Year() == date.Year()) && (v.DateTime.Month() == date.Month()) {
			ExpenseForTheMonth = append(ExpenseForTheMonth, v)
		}
	}
	return ExpenseForTheMonth
}

func getTotalForMonth(date time.Time, userID string) interface{} {
	var total float64
	for _, v := range Expensedb {
		if v.UserID == userID && (v.DateTime.Year() == date.Year()) && (v.DateTime.Month() == date.Month()) {
			total += v.Amount
		}
	}
	totalStruct := struct {
		Total float64 `json:"TotalExpenses"`
	}{
		total,
	}
	return totalStruct
}

func listExpenseBreakdownForMonth(date time.Time, userID string) *ExpenseBreakdown {
	var breakdown = ExpenseBreakdown{
		Savings: []*Expense{},
		Needs:   []*Expense{},
		Wants:   []*Expense{},
	}
	for _, v := range listExpensesForMonth(date, userID) {
		if v.Category == "Needs" {
			breakdown.Needs = append(breakdown.Needs, v)
		} else if v.Category == "Wants" {
			breakdown.Wants = append(breakdown.Wants, v)
		} else if v.Category == "Savings" {
			breakdown.Savings = append(breakdown.Savings, v)
		}
	}
	return &breakdown
}

func getExpenseBreakdownForMonth(date time.Time, userID string) interface{} {
	breakdown := listExpenseBreakdownForMonth(date, userID)
	var savingsTotal float64
	var needsTotal float64
	var wantsTotal float64

	for _, v := range breakdown.Savings {
		savingsTotal += v.Amount
	}
	for _, v := range breakdown.Needs {
		needsTotal += v.Amount
	}
	for _, v := range breakdown.Wants {
		wantsTotal += v.Amount
	}

	expenseBreakdown := struct {
		Savings float64
		Needs   float64
		Wants   float64
	}{
		savingsTotal,
		needsTotal,
		wantsTotal,
	}

	return expenseBreakdown
}

func validateUser(Email string, Password string) (bool, string) {
	for _, v := range Userdb {
		if v.Email == Email && v.Password == Password {
			return true, v.ID
		}
	}
	return false, ""
}

func extractClaims(tokenString string) (jwt.MapClaims, bool) {
	secretKey := []byte("secret-key")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}

func extractUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) < 2 {
		return "", false
	}
	token := strings.TrimSpace(splitToken[1])

	claims, ok := extractClaims(token)
	if !ok {
		return "", false
	}
	userID := fmt.Sprintf("%v", claims["UserID"])
	return userID, true
}
