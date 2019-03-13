# budget-tool


POST /budget
    adds a new budget to the list of budgets. Requires a name
GET /budget/{budgetName}/expenses
    returns a list of expenses for the selected budget
POST /budget/{budgetName}/expenses
    adds a new expense to the selected budget. Requires a name and price
GET /budget/budgetName/expenses/{expenseID}
    returns a single item
