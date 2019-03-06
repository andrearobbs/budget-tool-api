package budget

import (
	"database/sql"
	"fmt"
)

type Expense struct {
	Id       int
	Name     string
	Cost     float64
	BudgetId int
}

type Budget struct {
	Id   int
	Name string
}

type BudgetService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *BudgetService {
	return &BudgetService{
		db: db,
	}
}

const (
	insertNewBudgetQuery = "INSERT INTO budget (budget_name) VALUES (?)"

	selectLastId = "SELECT LAST_INSERT_ID()"

	selectBudgetQuery = "SELECT budget_id, budget_name FROM budget WHERE budget_name = ?"

	insertExpenseQuery = "INSERT INTO expense (expense_name, expense_cost, budget_id) VALUES (?,?,?)"

	selectExpensesQuery = "SELECT expense_id, expense_name, expense_cost, budget_id FROM expense WHERE budget_id = ?"
)

func (a *BudgetService) FindOrCreateBudget(budgetName string) (Budget, error) {
	trxn, err := a.db.Begin()
	if err != nil {
		return Budget{}, err
	}

	row := trxn.QueryRow(selectBudgetQuery, budgetName)

	var budget Budget

	err = row.Scan(
		&budget.Id,
		&budget.Name,
	)
	if err == nil {

		trxn.Commit()
		return budget, nil
	}
	fmt.Println(err)

	_, err = trxn.Exec(insertNewBudgetQuery, budgetName)
	if err != nil {
		trxn.Rollback()
		return Budget{}, err
	}

	row = trxn.QueryRow(selectLastId)

	err = row.Scan(
		&budget.Id,
	)
	if err != nil {
		trxn.Rollback()
		return Budget{}, err
	}

	budget.Name = budgetName

	trxn.Commit()

	return budget, nil
}

func (a *BudgetService) AddExpense(expense Expense) {

	_, err := a.db.Exec(insertExpenseQuery, expense.Name, expense.Cost, expense.BudgetId)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *BudgetService) ListExpenses(budgetId int) ([]Expense, error) {
	rows, err := a.db.Query(selectExpensesQuery, budgetId)
	if err != nil {
		return nil, err
	}

	var expenses []Expense
	for rows.Next() {
		var expense Expense

		err := rows.Scan(
			&expense.Id,
			&expense.Name,
			&expense.Cost,
			&expense.BudgetId,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (a *BudgetService) CalculateGrandTotal(expenses []Expense) float64 {
	var sum float64

	for _, x := range expenses {
		sum += x.Cost
	}

	return sum
}
