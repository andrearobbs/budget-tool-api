package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/andrearobbs/budget-tool/budget"
	_ "github.com/go-sql-driver/mysql"
	"github.com/manifoldco/promptui"
)

const (
	viewBudgetCmd = "View Budget"
	addExpenseCmd = "Add a New Expense"
)

type CLI struct {
	budgetService *budget.BudgetService
}

func New(budgetService *budget.BudgetService) *CLI {
	return &CLI{
		budgetService: budgetService,
	}
}

func (c *CLI) MainMenu() {

	nameBudgetPrompt := promptui.Prompt{
		Label: "Name a New Budget or Type the Name of an Existing Budget",
	}
	budgetname, err := nameBudgetPrompt.Run()
	if err != nil {
		return
	}

	theBudget, err := c.budgetService.FindOrCreateBudget(budgetname)
	if err != nil {
		fmt.Printf("Error Finding or Creating Budget %v\n", err)
		return
	}

	for {
		fmt.Println()

		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{
				viewBudgetCmd,
				addExpenseCmd,
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case viewBudgetCmd:
			c.ListExpenses(theBudget)

		case addExpenseCmd:
			err := c.AddExpense(theBudget)
			if err != nil {
				fmt.Printf("Prompt failed %n\n", err)
				return
			}
		}
	}

}

func (c *CLI) ListExpenses(theBudget budget.Budget) {
	budgetItems, err := c.budgetService.ListExpenses(theBudget.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, x := range budgetItems {
		fmt.Println(x.Name, x.Cost)
	}

	total := c.budgetService.CalculateGrandTotal(budgetItems)

	fmt.Println("Your grand total is $", total, ", mofo!")
}

func (c *CLI) AddExpense(theBudget budget.Budget) error {
	namePrompt := promptui.Prompt{
		Label: "Name this Expense",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	cost, err := c.numberPromptHelper("What is the total cost for this expense?")
	if err != nil {
		return err
	}

	newExpense := budget.Expense{
		Name:     name,
		Cost:     cost,
		BudgetId: theBudget.Id,
	}

	c.budgetService.AddExpense(newExpense)

	fmt.Println("Added new expense to budget!")

	return nil
}

func (c *CLI) numberPromptHelper(label string) (float64, error) {

	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("You need to type a number, dummy!")
		}
		return nil
	}

	costPrompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	numberStr, err := costPrompt.Run()
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}
