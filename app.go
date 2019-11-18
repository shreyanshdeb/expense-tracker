package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

func getExpenses(userID string) ([]Expense, bool) {
	var result []Expense
	iter := getAllExpensesFromDb(userID)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, false
		}
		var expense Expense
		mapstructure.Decode(doc.Data(), &expense)
		result = append(result, expense)
	}
	return result, true
}

func addExpense(expense Expense) bool {
	ok := addExpenseToDb(expense)
	return ok
}

func getExpenseByID(userID string, ID string) (Expense, bool) {
	iter := getExpenseByIDFromDb(userID, ID)
	var expense Expense
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return expense, false
		}
		mapstructure.Decode(doc.Data(), &expense)
	}
	return expense, true
}

func deleteExpenseByID(userID string, ID string) bool {
	ok := deleteExpenseFromDb(userID, ID)
	return ok
}

func updateExpenseByID(id string, expense Expense) bool {
	ok := updateExpenseInDb(id, expense)
	return ok
}

func listExpensesForMonth(userID string, fromDate time.Time, toDate time.Time) ([]*Expense, bool) {
	// var ExpenseForTheMonth = []*Expense{}
	// for _, v := range Expensedb {
	// 	if v.UserID == userID && (v.DateTime.Year() == date.Year()) && (v.DateTime.Month() == date.Month()) {
	// 		ExpenseForTheMonth = append(ExpenseForTheMonth, v)
	// 	}
	// }
	// return ExpenseForTheMonth
	iter := getExpensesForDateFromDb(userID, fromDate, toDate)
	var result []*Expense
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, false
		}
		var expense Expense
		err = mapstructure.Decode(doc.Data(), &expense)
		if err != nil {
			return nil, false
		}
		result = append(result, &expense)
	}
	return result, true
}

func getTotalForMonth(userID string, fromDate time.Time, toDate time.Time) (interface{}, bool) {
	var total float64
	expenseList, ok := listExpensesForMonth(userID, fromDate, toDate)
	if !ok {
		return nil, false
	}
	for _, v := range expenseList {
		total += v.Amount
	}
	totalStruct := struct {
		Total float64 `json:"TotalExpenses"`
	}{
		total,
	}
	return totalStruct, true
}

func listExpenseBreakdownForMonth(userID string, fromDate time.Time, toDate time.Time) (*ExpenseBreakdown, bool) {
	var breakdown = ExpenseBreakdown{
		Savings: []*Expense{},
		Needs:   []*Expense{},
		Wants:   []*Expense{},
	}
	expenseList, ok := listExpensesForMonth(userID, fromDate, toDate)
	if !ok {
		return nil, false
	}
	for _, v := range expenseList {
		if v.Category == "Needs" {
			breakdown.Needs = append(breakdown.Needs, v)
		} else if v.Category == "Wants" {
			breakdown.Wants = append(breakdown.Wants, v)
		} else if v.Category == "Savings" {
			breakdown.Savings = append(breakdown.Savings, v)
		}
	}
	return &breakdown, true
}

func getExpenseBreakdownForMonth(userID string, fromDate time.Time, toDate time.Time) (interface{}, bool) {
	breakdown, ok := listExpenseBreakdownForMonth(userID, fromDate, toDate)
	if !ok {
		return nil, false
	}
	var savingsTotal float64
	var needsTotal float64
	var wantsTotal float64
	var total float64

	for _, v := range breakdown.Savings {
		savingsTotal += v.Amount
	}
	for _, v := range breakdown.Needs {
		needsTotal += v.Amount
	}
	for _, v := range breakdown.Wants {
		wantsTotal += v.Amount
	}
	total = savingsTotal + needsTotal + wantsTotal
	expenseBreakdown := struct {
		Total   float64
		Savings float64
		Needs   float64
		Wants   float64
	}{
		total,
		savingsTotal,
		needsTotal,
		wantsTotal,
	}

	return expenseBreakdown, true
}

func validateUser(Email string, Password string) (bool, string) {
	iter := getUserByEmail(Email)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, ""
		}
		var user User
		err = mapstructure.Decode(doc.Data(), &user)
		if err != nil {
			return false, ""
		}
		if comparePasswords(user.Password, Password) {
			return true, user.ID
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

func getTerminalDates(date time.Time) (time.Time, time.Time) {
	Year, Month, _ := date.Date()
	firstOfMonth := time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := time.Date(Year, Month+1, 1, 0, 0, 0, -1, time.UTC)
	return firstOfMonth, lastOfMonth
}

func signUp(user User) bool {
	iter := getUserByEmail(user.Email)
	var existingUser User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false
		}
		err = mapstructure.Decode(doc.Data(), existingUser)
		if err != nil {
			return false
		}
		break
	}
	if existingUser.ID == "" {
		id, _ := uuid.NewV4()
		user.ID = id.String()
		user.Password = hashAndSalt(user.Password)
		ok := addUserToDb(user)
		if !ok {
			return false
		}
		return true
	}
	return false
}

func signIn(user User) (bool, string, string) {
	ok, userid := validateUser(user.Email, user.Password)
	if !ok {
		return false, "", ""
	}
	var token string
	ok, token = generateToken(userid, user.Email)
	if !ok {
		return false, "", ""
	}
	return true, userid, token
}

func generateToken(userID string, email string) (bool, string) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return false, ""
	}
	return true, tokenString
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}
