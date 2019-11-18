package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	// "golang.org/x/net/context"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func getExpensesHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var expenseList []Expense
	expenseList, ok = getExpenses(userID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var JSONExpenseList []JSONExpense
	for _, v := range expenseList {
		JSONExpenseList = append(JSONExpenseList, *NewJSONExpense(&v))
	}
	jsonData, _ := json.Marshal(JSONExpenseList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]

	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var expense Expense
	expense, ok = getExpenseByID(userID, ID)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if expense.ID == "" {
		jsonData, _ := json.Marshal(ResponseStruct{
			false,
			"Did not find expense for ID: " + ID,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}
	JSONResult := NewJSONExpense(&expense)
	jsonData, _ := json.Marshal(JSONResult)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func addExpenseHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var je JSONExpense
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()

	id, _ := uuid.NewV4()
	expense.ID = id.String()
	expense.UserID = userID
	ok = addExpense(*expense)
	var jsonData []byte
	if !ok {
		jsonData, _ = json.Marshal(ResponseStruct{
			false,
			nil,
		})
	} else {
		jsonData, _ = json.Marshal(ResponseStruct{
			true,
			*expense,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	deleted := deleteExpenseByID(userID, id)
	if !deleted {
		jsonData, _ := json.Marshal(ResponseStruct{
			false,
			nil,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}
	jsonData, _ := json.Marshal(ResponseStruct{
		true,
		nil,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func updateExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var je JSONExpense

	params := mux.Vars(r)
	id := params["id"]
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()

	//updated := updateExpenseByID(id, *expense, userID)
	updated := updateExpenseByID(id, *expense)

	if !updated {
		jsonData, _ := json.Marshal(ResponseStruct{
			false,
			nil,
		})
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}

	jsonData, _ := json.Marshal(ResponseStruct{
		true,
		nil,
	})
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)
}

func listExpensesForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	date := time.Unix(monthInt, 0)
	formdate, todate := getTerminalDates(date)
	var expenseList []*Expense
	expenseList, ok = listExpensesForMonth(userID, formdate, todate)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var JSONExpenseList []JSONExpense
	for _, v := range expenseList {
		JSONExpenseList = append(JSONExpenseList, *NewJSONExpense(v))
	}
	jsonData, _ := json.Marshal(JSONExpenseList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getTotalForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		fmt.Println("Failed to convert time to int64")
	}
	var total interface{}
	fromDate, toDate := getTerminalDates(time.Unix(monthInt, 0))
	total, ok = getTotalForMonth(userID, fromDate, toDate)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(total)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func listExpenseBreakdownForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		fmt.Println("Failed to convert time to int64")
	}
	fromDate, toDate := getTerminalDates(time.Unix(monthInt, 0))
	breakdown, ok := listExpenseBreakdownForMonth(userID, fromDate, toDate)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(breakdown)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getExpenseBreakdownForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		fmt.Println("Failed to convert time to int64")
	}
	fromDate, toDate := getTerminalDates(time.Unix(monthInt, 0))
	breakdown, ok := getExpenseBreakdownForMonth(userID, fromDate, toDate)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(breakdown)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok, userid, token := signIn(user)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	response := struct {
		ID            string
		Email         string
		Authorization string
	}{
		userid,
		user.Email,
		token,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func isAuthorized(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer")
		token := strings.TrimSpace(splitToken[1])
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func testHandler(w http.ResponseWriter, r *http.Request) {
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok := signUp(user)
	var response ResponseStruct
	if !ok {
		response = ResponseStruct{
			false,
			nil,
		}

	} else {
		response = ResponseStruct{
			true,
			nil,
		}
	}
	responseJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	return
}
