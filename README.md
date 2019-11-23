#Expense Tracker

##Introduction

Expense tracker that categorizes expenses into three types 'Savings', 'Needs' and 'Wants'. Borrowing the idea from the famous budgeting rule 50/30/20. 

>This is a popular rule for breaking down your budget. The 50-30-20 rule puts 50% of your income toward necessities, like housing and bills. Twenty percent should then go toward financial goals, like paying off debt or saving for retirement. Finally, 30% of your income can be allocated to wants, like dining or entertainment.

>-Lifehacker, The 10 Best Financial Rules of Thumb

##Behind the scenes

* The webserver is written in go with routing in gorilla mux. 

* Data is stored in Cloud Firestore.

* Auth using JSON Web Tokens.

##What did I Learn

* Creating Restful APIs in go.

* Gracefully handling errors in go.

* Creating and reading from a config file in go.

* Connecting and storing data in Cloud Firestore.
