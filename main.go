package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/andrearobbs/budget-tool-api/budget"
	"github.com/andrearobbs/budget-tool-api/db"
	"github.com/andrearobbs/budget-tool-api/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const port = ":8000"

func main() {

	db, err := db.ConnectDatabase("budget_app_db.config")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	budgetService := budget.NewService(db)

	budgetServer := server.NewServer(budgetService)

	router := mux.NewRouter()
	router.HandleFunc("/budget", budgetServer.CreateBudgetHandler).Methods("POST")
	router.HandleFunc("/budget/{budgetName}/expenses", budgetServer.ListExpensesHandler).Methods("GET")
	//router.HandleFunc("/budget/{budgetName}/expenses", budgetServer.AddExpenseToBudgetHandler).Methods("POST")
	//router.HandleFunc("/budget/budgetName/expenses/{expenseID}", budgetServer.GetSingleExpenseHandler).Methods("GET")

	http.Handle("/", router)

	fmt.Println("Waiting for requests on port:", port)
	http.ListenAndServe(port, nil)
}
