package main

import (
	"time"
)

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

func getExpenseByID(id string) (*Expense, error) {
	for _, v := range Expensedb {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, ErrExpenseNotFound
}

func deleteExpenseByID(id string) bool {
	for i, v := range Expensedb {
		if v.ID == id {
			Expensedb = append(Expensedb[:i], Expensedb[i+1:]...)
			return true
		}
	}
	return false
}

func updateExpenseByID(id string, expense Expense) bool {
	for _, v := range Expensedb {
		if v.ID == id {
			v.Title = expense.Title
			v.Category = expense.Category
			v.Amount = expense.Amount
			v.DateTime = expense.DateTime
			v.UserID = expense.UserID
			return true
		}
	}
	return false
}

func getExpensesForCurrentMonth() []*Expense {
	var ExpenseForTheMonth = []*Expense{}
	currentTime := time.Now()
	for _, v := range Expensedb {
		if (v.DateTime.Year() == currentTime.Year()) && (v.DateTime.Month() == currentTime.Month()) {
			ExpenseForTheMonth = append(ExpenseForTheMonth, v)
		}
	}
	return ExpenseForTheMonth
}

func getTotalForCurrentMonth() interface{} {
	var total float64
	currentTime := time.Now()
	for _, v := range Expensedb {
		if (v.DateTime.Year() == currentTime.Year()) && (v.DateTime.Month() == currentTime.Month()) {
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

func listExpenseBreakdownForCurrentMonth() *ExpenseBreakdown {
	var breakdown = ExpenseBreakdown{
		Savings: []*Expense{},
		Needs:   []*Expense{},
		Wants:   []*Expense{},
	}
	for _, v := range getExpensesForCurrentMonth() {
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

func getExpenseBreakdownForCurrentMonth() interface{} {
	breakdown := listExpenseBreakdownForCurrentMonth()
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

//------------------------------------------------------------------------------

func listExpensesForMonth(date time.Time) []*Expense {
	var ExpenseForTheMonth = []*Expense{}
	for _, v := range Expensedb {
		if (v.DateTime.Year() == date.Year()) && (v.DateTime.Month() == date.Month()) {
			ExpenseForTheMonth = append(ExpenseForTheMonth, v)
		}
	}
	return ExpenseForTheMonth
}

func getTotalForMonth(date time.Time) interface{} {
	var total float64
	for _, v := range Expensedb {
		if (v.DateTime.Year() == date.Year()) && (v.DateTime.Month() == date.Month()) {
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

func listExpenseBreakdownForMonth(date time.Time) *ExpenseBreakdown {
	var breakdown = ExpenseBreakdown{
		Savings: []*Expense{},
		Needs:   []*Expense{},
		Wants:   []*Expense{},
	}
	for _, v := range listExpensesForMonth(date) {
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

func getExpenseBreakdownForMonth(date time.Time) interface{} {
	breakdown := listExpenseBreakdownForMonth(date)
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
