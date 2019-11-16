package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	firestore "cloud.google.com/go/firestore"
)

func initLocalDb() {
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

var firestoreClient *firestore.Client
var config Configuration

func readConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	readConfig()
	initLocalDb()

	var err error
	firestoreClient, err = initDb()
	if err != nil {
		log.Fatalln("Error initialization Firestore", err)
	}
	//defer firestoreClient.Close()

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

	r.HandleFunc("/test", testHandler).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	fmt.Println("Listening on ", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
