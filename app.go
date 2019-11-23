package main

import (
	"encoding/json"
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

func getExpenses(userID string) ([]Expense, error) {
	var result []Expense
	iter := getAllExpensesFromDb(userID)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var expense Expense
		mapstructure.Decode(doc.Data(), &expense)
		result = append(result, expense)
	}
	return result, nil
}

func addExpense(expense Expense) error {
	err := addExpenseToDb(expense)
	if err != nil {
		return err
	}
	return nil
}

func getExpenseByID(userID string, ID string) (Expense, error) {
	iter := getExpenseByIDFromDb(userID, ID)
	var expense Expense
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return expense, err
		}
		mapstructure.Decode(doc.Data(), &expense)
	}
	return expense, nil
}

func deleteExpenseByID(userID string, ID string) error {
	err := deleteExpenseFromDb(userID, ID)
	if err != nil {
		return err
	}
	return nil
}

func updateExpenseByID(id string, expense Expense) error {
	err := updateExpenseInDb(id, expense)
	if err != nil {
		return err
	}
	return nil
}

func listExpensesForMonth(userID string, fromDate time.Time, toDate time.Time) ([]*Expense, error) {
	iter := getExpensesForDateFromDb(userID, fromDate, toDate)
	var result []*Expense
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var expense Expense
		err = mapstructure.Decode(doc.Data(), &expense)
		if err != nil {
			return nil, err
		}
		result = append(result, &expense)
	}
	return result, nil
}

func getTotalForMonth(userID string, fromDate time.Time, toDate time.Time) (interface{}, error) {
	var total float64
	expenseList, err := listExpensesForMonth(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	for _, v := range expenseList {
		total += v.Amount
	}
	totalStruct := struct {
		Total float64 `json:"TotalExpenses"`
	}{
		total,
	}
	return totalStruct, nil
}

func listExpenseBreakdownForMonth(userID string, fromDate time.Time, toDate time.Time) (*ExpenseBreakdown, error) {
	var breakdown = ExpenseBreakdown{
		Savings: []*Expense{},
		Needs:   []*Expense{},
		Wants:   []*Expense{},
	}
	expenseList, err := listExpensesForMonth(userID, fromDate, toDate)
	if err != nil {
		return nil, err
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
	return &breakdown, nil
}

func getExpenseBreakdownForMonth(userID string, fromDate time.Time, toDate time.Time) (interface{}, error) {
	breakdown, err := listExpenseBreakdownForMonth(userID, fromDate, toDate)
	if err != nil {
		return nil, err
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

	return expenseBreakdown, nil
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

func extractClaims(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte("secret-key")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf(err.Error(), "Could not extract claim")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Invalid Claims")
	}
}

func extractUserID(w http.ResponseWriter, r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) < 2 {
		return "", fmt.Errorf("Invalid bearer token")
	}
	token := strings.TrimSpace(splitToken[1])

	claims, err := extractClaims(token)
	if err != nil {
		return "", err
	}
	userID := fmt.Sprintf("%v", claims["UserID"])
	return userID, nil
}

func getTerminalDates(date time.Time) (time.Time, time.Time) {
	Year, Month, _ := date.Date()
	firstOfMonth := time.Date(Year, Month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := time.Date(Year, Month+1, 1, 0, 0, 0, -1, time.UTC)
	return firstOfMonth, lastOfMonth
}

func signUp(user User) error {
	iter := getUserByEmail(user.Email)
	var existingUser User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		err = mapstructure.Decode(doc.Data(), existingUser)
		if err != nil {
			return err
		}
		break
	}
	if existingUser.ID != "" {
		return fmt.Errorf("User already exists")
	}
	id, _ := uuid.NewV4()
	user.ID = id.String()
	user.Password = hashAndSalt(user.Password)
	err := addUserToDb(user)
	if err != nil {
		return err
	}
	return nil
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

func createProblemJSON(Type string, Title string, Status string, Detail string, Instance string) []byte {
	problemJSON, err := json.Marshal(ProblemDetails{
		Type,
		Title,
		Status,
		Detail,
		Instance,
	})
	fmt.Println(err)
	return problemJSON
}
