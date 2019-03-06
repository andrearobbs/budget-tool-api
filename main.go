package main

import (
	"fmt"
	"os"

	"github.com/andrearobbs/budget-tool/budget"
	"github.com/andrearobbs/budget-tool/cli"
	"github.com/andrearobbs/budget-tool/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.ConnectDatabase("budget_app_db.config")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	budgetService := budget.NewService(db)

	cliMenu := cli.New(budgetService)

	cliMenu.MainMenu()

}
