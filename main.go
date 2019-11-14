package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func initDb() {
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	Expensedb = append(Expensedb, &Expense{
		id1.String(),
		"Rent",
		"Needs",
		4500,
		time.Unix(1572609600, 0),
		"ac8421b4-4d05-4b1b-881f-90a7a7e0108b",
	})
	Expensedb = append(Expensedb, &Expense{
		id2.String(),
		"House Help",
		"Needs",
		5000,
		time.Unix(1572991215, 0),
		"ac8421b4-4d05-4b1b-881f-90a7a7e0108b",
	})
	Expensedb = append(Expensedb, &Expense{
		id3.String(),
		"House Help",
		"Needs",
		5000,
		time.Unix(1572991215, 0),
		"09c9d7bf-de35-41a3-a32f-f18bf9fe5c94",
	})
	Userdb = append(Userdb, &User{
		"ac8421b4-4d05-4b1b-881f-90a7a7e0108b",
		"Tyler",
		"tylerj@top.com",
		"9713161575",
		"dema",
	})
	Userdb = append(Userdb, &User{
		"09c9d7bf-de35-41a3-a32f-f18bf9fe5c94",
		"Josh",
		"joshd@top.com",
		"9826342221",
		"scottland",
	})
}

func main() {
	initDb()

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/expenses", isAuthorized(getExpensesHandler)).Methods("GET")
	r.HandleFunc("/api/v1/expense/{id}", isAuthorized(getExpenseByIDHandler)).Methods("GET")
	r.HandleFunc("/api/v1/expense", isAuthorized(addExpenseHandler)).Methods("POST")
	r.HandleFunc("/api/v1/expense/{id}", isAuthorized(deleteExpenseByIDHandler)).Methods("DELETE")
	r.HandleFunc("/api/v1/expense/{id}", isAuthorized(updateExpenseByIDHandler)).Methods("PUT")

	r.HandleFunc("/api/v1/expensesformonth/{month}", isAuthorized(listExpensesForMonthHandler)).Methods("GET")
	r.HandleFunc("/api/v1/totalexpenseformonth/{month}", isAuthorized(getTotalForMonthHandler)).Methods("GET")
	r.HandleFunc("/api/v1/breakdownformonth/{month}", isAuthorized(listExpenseBreakdownForMonthHandler)).Methods("GET")
	r.HandleFunc("/api/v1/sumbreakdownformonth/{month}", isAuthorized(getExpenseBreakdownForMonthHandler)).Methods("GET")

	r.HandleFunc("/api/v1/signin", signIn).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	fmt.Println("Listening on 1234")
	log.Fatal(http.ListenAndServe(":1234", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
