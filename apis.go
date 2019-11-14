package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func getExpensesHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(getExpenses(userID))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expense, err := getExpenseByID(id, userID)
	if err != nil {
		jsonData, _ := json.Marshal(ResponseStruct{
			false,
			nil,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}
	jsonData, _ := json.Marshal(expense)
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
	addExpense(*expense)

	jsonData, _ := json.Marshal(ResponseStruct{
		true,
		expense,
	})
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

	deleted := deleteExpenseByID(id, userID)
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
	userID, ok := extractUserID(w, r)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var je JSONExpense

	params := mux.Vars(r)
	id := params["id"]
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()

	updated := updateExpenseByID(id, *expense, userID)

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
		fmt.Println("Failed to convert time to int64")
	}
	jsonData, _ := json.Marshal(listExpensesForMonth(time.Unix(monthInt, 0), userID))
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
	total := getTotalForMonth(time.Unix(monthInt, 0), userID)
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
	jsonData, _ := json.Marshal(listExpenseBreakdownForMonth(time.Unix(monthInt, 0), userID))
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
	jsonData, _ := json.Marshal(getExpenseBreakdownForMonth(time.Unix(monthInt, 0), userID))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ok, userid := validateUser(user.Email, user.Password)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userid,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := struct {
		ID             string
		Email          string
		Authorization  string
		ExpirationTime string
	}{
		userid,
		user.Email,
		tokenString,
		expirationTime.String(),
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
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
