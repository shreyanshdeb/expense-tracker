package main

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
			v.Category = expense.Category
			v.Amount = expense.Amount
			v.DateTime = expense.DateTime
			v.UserID = expense.UserID
			return true
		}
	}
	return false
}
