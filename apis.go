package main

import (
	"encoding/json"
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
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	var expenseList []Expense
	expenseList, err = getExpenses(userID)
	if err != nil {
		problemJSON := createProblemJSON(ExpensesProblemURI, "Could not fetch expenses", "InternalServerError", "Problem in fetching expenses", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
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

	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	expense, err := getExpenseByID(userID, ID)
	if err != nil || expense.ID == "" {
		var errorString string
		if err == nil {
			errorString = "Could not find expense with the provided ID"
		} else {
			errorString = err.Error()
		}
		problemJSON := createProblemJSON(ExpensesProblemURI, "Could not fetch expenses", "InternalServerError", "Problem in fetching expenses", errorString)
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	JSONResult := NewJSONExpense(&expense)
	jsonData, _ := json.Marshal(JSONResult)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func addExpenseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	var je JSONExpense
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()

	id, _ := uuid.NewV4()
	expense.ID = id.String()
	expense.UserID = userID
	err = addExpense(*expense)

	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Could not add expenses", "InternalServerError", "Problem in adding expense", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}

	var jsonData []byte
	jsonData, _ = json.Marshal(ResponseStruct{
		true,
		*expense,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func deleteExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	err = deleteExpenseByID(userID, id)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Could not delete expenses", "InternalServerError", "Problem in deleting expense", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
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
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}

	var je JSONExpense

	params := mux.Vars(r)
	id := params["id"]
	_ = json.NewDecoder(r.Body).Decode(&je)

	expense := je.Expense()
	expense.ID = id
	expense.UserID = userID
	err = updateExpenseByID(id, *expense)

	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Could not update the expenses", "InternalServerError", "Problem in updating the expense", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
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
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		problemJSON := createProblemJSON(ParsingProblemURL, "Could not parse the date", "InternalServerError", "Problem in parsing the date", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	date := time.Unix(monthInt, 0)
	formdate, todate := getTerminalDates(date)
	var expenseList []*Expense
	expenseList, err = listExpensesForMonth(userID, formdate, todate)
	if err != nil {
		problemJSON := createProblemJSON(ExpensesProblemURI, "Could not fetch expenses", "InternalServerError", "Problem in fetching expenses", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
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
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		problemJSON := createProblemJSON(ParsingProblemURL, "Could not parse the date", "InternalServerError", "Problem in parsing the date", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	var total interface{}
	fromDate, toDate := getTerminalDates(time.Unix(monthInt, 0))
	total, err = getTotalForMonth(userID, fromDate, toDate)
	if err != nil {
		problemJSON := createProblemJSON(ExpensesProblemURI, "Could not fetch expenses", "InternalServerError", "Problem in fetching expenses", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	jsonData, _ := json.Marshal(total)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func listExpenseBreakdownForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		problemJSON := createProblemJSON(ParsingProblemURL, "Could not parse the date", "InternalServerError", "Problem in parsing the date", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	fromDate, toDate := getTerminalDates(time.Unix(monthInt, 0))
	breakdown, err := listExpenseBreakdownForMonth(userID, fromDate, toDate)
	if err != nil {
		problemJSON := createProblemJSON(ExpensesProblemURI, "Could not fetch expenses", "InternalServerError", "Problem in fetching expenses", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	jsonData, _ := json.Marshal(breakdown)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getExpenseBreakdownForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserID(w, r)
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
		return
	}
	params := mux.Vars(r)
	month := params["month"]
	monthInt, err := strconv.ParseInt(month, 10, 64)
	if err != nil {
		problemJSON := createProblemJSON(ParsingProblemURL, "Could not parse the date", "InternalServerError", "Problem in parsing the date", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
		return
	}
	fromDate, toDate := getTerminalDates(time.Unix(monthInt, 0))
	breakdown, err := getExpenseBreakdownForMonth(userID, fromDate, toDate)
	if err != nil {
		problemJSON := createProblemJSON(ExpensesProblemURI, "Could not fetch expenses", "InternalServerError", "Problem in fetching expenses", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(problemJSON)
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
		problemJSON := createProblemJSON(InvalidTokenURI, "Incorrect Request Body", "StatusBadRequest", "Problem with Request Body", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(problemJSON)
		return
	}
	ok, userid, token := signIn(user)
	if !ok {
		problemJSON := createProblemJSON(InvalidTokenURI, "Invalid Username or Password", "Unauthorized", "Problem with Username or Password", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(problemJSON)
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
		w.WriteHeader(http.StatusInternalServerError)
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
			problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", "")
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(problemJSON)
			return
		}
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})
		if err != nil {
			problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(problemJSON)
			return
		}
		if !tkn.Valid {
			problemJSON := createProblemJSON(InvalidTokenURI, "Invalid auth token", "Unauthorized", "Problem with the auth token", err.Error())
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(problemJSON)
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
		problemJSON := createProblemJSON(InvalidTokenURI, "Incorrect Request Body", "StatusBadRequest", "Problem with Request Body", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(problemJSON)
		return
	}
	err = signUp(user)
	var response ResponseStruct
	if err != nil {
		problemJSON := createProblemJSON(InvalidTokenURI, "Cannot signup", "Unauthorized", "Problem in signing up", err.Error())
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(problemJSON)
		return
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
