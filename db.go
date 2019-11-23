package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

func initDb() (*firestore.Client, error) {
	opt := option.WithCredentialsFile(config.FirebaseCredentialsFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
		return nil, err
	}
	fmt.Println("Connected to Firebase")

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to Firestore")
	return client, nil
}

func addExpenseToDb(expense Expense) error {
	_, err := firestoreClient.Collection(config.Collections["Expense"]).Doc(expense.ID).Set(context.Background(), expense)
	if err != nil {
		return fmt.Errorf(err.Error(), "Could not create expense")
	}
	return nil
}

func getAllExpensesFromDb(userID string) *firestore.DocumentIterator {
	iter := firestoreClient.Collection(config.Collections["Expense"]).Where("UserID", "==", userID).Documents(context.Background())
	return iter
}

func getExpenseByIDFromDb(userID string, ID string) *firestore.DocumentIterator {
	iter := firestoreClient.Collection(config.Collections["Expense"]).Where("UserID", "==", userID).Where("ID", "==", ID).Documents(context.Background())
	return iter
}

func getExpensesForDateFromDb(userID string, fromdate time.Time, todate time.Time) *firestore.DocumentIterator {
	iter := firestoreClient.Collection(config.Collections["Expense"]).Where("UserID", "==", userID).Where("DateTime", ">=", fromdate).Where("DateTime", "<=", todate).Documents(context.Background())
	return iter
}

func deleteExpenseFromDb(userID string, ID string) error {
	_, err := firestoreClient.Collection(config.Collections["Expense"]).Doc(ID).Delete(context.Background())
	if err != nil {
		return fmt.Errorf(err.Error(), "Could not delete expense")
	}
	return nil
}

func updateExpenseInDb(ID string, expense Expense) error {
	_, err := firestoreClient.Collection(config.Collections["Expense"]).Doc(ID).Set(context.Background(), expense)
	if err != nil {
		return fmt.Errorf(err.Error(), "Could not update expense")
	}
	return nil
}

func addUserToDb(user User) error {
	_, err := firestoreClient.Collection(config.Collections["User"]).Doc(user.ID).Set(context.Background(), user)
	if err != nil {
		return fmt.Errorf(err.Error(), "Could not create user")
	}
	return nil
}

func getUserByID(userID string) *firestore.DocumentIterator {
	iter := firestoreClient.Collection(config.Collections["User"]).Where("ID", "==", userID).Documents(context.Background())
	return iter
}

func getUserByEmail(userName string) *firestore.DocumentIterator {
	iter := firestoreClient.Collection(config.Collections["User"]).Where("Email", "==", userName).Documents(context.Background())
	return iter
}

func getUserByEmailAndPassword(userName string, password string) *firestore.DocumentIterator {
	iter := firestoreClient.Collection(config.Collections["User"]).Where("Email", "==", userName).Where("Password", "==", password).Documents(context.Background())
	return iter
}
