package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andrearobbs/budget-tool-api/budget"
	"github.com/gorilla/mux"
)

type BudgetServer struct {
	budgetService *budget.BudgetService
}

func NewServer(budgetService *budget.BudgetService) *BudgetServer {
	return &BudgetServer{
		budgetService: budgetService,
	}
}

type CreateBudgetRequest struct {
	Name string
}

func (s *BudgetServer) CreateBudgetHandler(rw http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var newBudget CreateBudgetRequest
	err = json.Unmarshal(requestBody, &newBudget)
	if err != nil {
		fmt.Println("Error unmarshaling new budget details:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.budgetService.FindOrCreateBudget(newBudget.Name)
	if err != nil {
		fmt.Println("Error creating new budget:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (s *BudgetServer) ListExpensesHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	budgetNameStr := vars["budgetName"]

	budget, err := s.budgetService.FindOrCreateBudget(budgetNameStr)
	if err != nil {
		fmt.Println("error finding or creating budget:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	expenses, err := s.budgetService.ListExpenses(budget.Id)
	if err != nil {
		fmt.Println("Error listing expenses:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	expensesJSON, err := json.Marshal(expenses)
	if err != nil {
		fmt.Println("Error marshaling expenses:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.Write(expensesJSON)
}

type AddExpenseToBudgetRequest struct {
	ExpenseName  string
	ExpensePrice int
}

// func (s *BudgetServer) AddExpenseToBudgetHandler(rw http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	expenseNameStr := vars["budgetName"]

// 	requestBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Println("Error reading request body:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	var newExpenseDetails AddExpenseToBudgetRequest
// 	err = json.Unmarshal(requestBody, &newGameDetails)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling game cabinet details:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	err = s.budgetService.AddExpenseToBudget(newExpenseDetails.ExpenseName, budgetName, newExpenseDetails.ExpensePrice)
// 	if err != nil {
// 		fmt.Println("Error adding expense to budget:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	rw.WriteHeader(http.StatusCreated)
// }

// func (s *BudgetServer) GetSingleExpenseHandler(rw http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	expenseNameStr := vars["expenseName"]

// 	expenseName, err := strconv.Atoi(expenseNameStr)
// 	if err != nil {
// 		fmt.Println("Invalid expenseName:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	expense, err := s.budgetService.GetSingleExpenseHandler(expenseName)
// 	if err != nil {
// 		fmt.Println("Error getting expense:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	expenseJSON, err := json.Marshal(expense)
// 	if err != nil {
// 		fmt.Println("Error marshaling expense:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	rw.Header().Add("Content-Type", "application/json")
// 	rw.Write(expenseJSON)
// }
