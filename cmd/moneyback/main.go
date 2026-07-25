package main

import (
	"flag"
	"fmt"
	"os"

	"money-backword/internal/api"
	"money-backword/internal/ledger"
	"money-backword/internal/storage"
)

var (
	dbPath   = flag.String("db", "./money.json", "path to database file")
	apiAddr  = flag.String("api", ":8080", "HTTP API listen address")
	verboseF = flag.Bool("v", false, "verbose output")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(0)
	}

	cmd := flag.Arg(0)

	store, err := storage.NewJSONStore(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening database: %v\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "add":
		cmdAdd(flag.Args()[1:], store)
	case "list":
		cmdList(flag.Args()[1:], store)
	case "accounts":
		cmdAccounts(flag.Args()[1:], store)
	case "report":
		cmdReport(flag.Args()[1:], store)
	case "budget":
		cmdBudget(flag.Args()[1:], store)
	case "serve":
		cmdServe(store)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func cmdAdd(args []string, store *storage.JSONStore) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	accountID := fs.String("account", "", "account ID")
	category := fs.String("category", "uncategorized", "transaction category")
	amount := fs.Float64("amount", 0, "transaction amount")
	description := fs.String("desc", "", "transaction description")

	fs.Parse(args)

	if *accountID == "" || *amount == 0 {
		fmt.Fprintf(os.Stderr, "missing required flags: -account and -amount\n")
		return
	}

	txn := ledger.NewTransaction(*accountID, *category, *amount, *description)
	if err := store.AddTransaction(txn); err != nil {
		fmt.Fprintf(os.Stderr, "error adding transaction: %v\n", err)
		return
	}

	fmt.Printf("added transaction %s\n", txn.ID)
}

func cmdList(args []string, store *storage.JSONStore) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	accountID := fs.String("account", "", "filter by account ID")
	limit := fs.Int("limit", 20, "number of transactions to show")

	fs.Parse(args)

	txns, err := store.GetTransactions(*accountID, *limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching transactions: %v\n", err)
		return
	}

	if len(txns) == 0 {
		fmt.Println("no transactions found")
		return
	}

	fmt.Printf("%-36s %-12s %-15s %-20s %s\n", "ID", "Amount", "Category", "Account", "Description")
	fmt.Println(repeatStr("-", 100))
	for _, txn := range txns {
		fmt.Printf("%-36s %-12.2f %-15s %-20s %s\n", txn.ID, txn.Amount, txn.Category, txn.AccountID, txn.Description)
	}
}

func cmdAccounts(args []string, store *storage.JSONStore) {
	accounts, err := store.GetAllAccounts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching accounts: %v\n", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Println("no accounts configured")
		return
	}

	fmt.Printf("%-20s %-15s %s\n", "ID", "Balance", "Name")
	fmt.Println(repeatStr("-", 60))
	for _, acc := range accounts {
		balance := acc.GetBalance()
		fmt.Printf("%-20s %-15.2f %s\n", acc.ID, balance, acc.Name)
	}
}

func cmdReport(args []string, store *storage.JSONStore) {
	fs := flag.NewFlagSet("report", flag.ExitOnError)
	monthStr := fs.String("month", "", "month (YYYY-MM)")
	reportType := fs.String("type", "summary", "report type: summary, category, csv")

	fs.Parse(args)

	// TODO: implement different report types and month filtering
	report, err := store.GenerateMonthlyReport(*monthStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating report: %v\n", err)
		return
	}

	fmt.Printf("Total Income: %.2f\n", report.TotalIncome)
	fmt.Printf("Total Expense: %.2f\n", report.TotalExpense)
	fmt.Printf("Net: %.2f\n", report.Net)

	if *reportType == "category" && len(report.CategoryBreakdown) > 0 {
		fmt.Println("\nBy Category:")
		for cat, amount := range report.CategoryBreakdown {
			fmt.Printf("  %s: %.2f\n", cat, amount)
		}
	}
}

func cmdBudget(args []string, store *storage.JSONStore) {
	fs := flag.NewFlagSet("budget", flag.ExitOnError)
	category := fs.String("category", "", "category name")
	limit := fs.Float64("limit", 0, "monthly budget limit")
	action := fs.String("action", "list", "action: list, set, check")

	fs.Parse(args)

	switch *action {
	case "set":
		if *category == "" || *limit == 0 {
			fmt.Fprintf(os.Stderr, "budget set: missing -category and -limit\n")
			return
		}
		budget := ledger.NewBudget(*category, *limit)
		if err := store.SetBudget(budget); err != nil {
			fmt.Fprintf(os.Stderr, "error setting budget: %v\n", err)
			return
		}
		fmt.Printf("set budget for %s: %.2f\n", *category, *limit)

	case "list", "check":
		budgets, err := store.GetAllBudgets()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error fetching budgets: %v\n", err)
			return
		}
		if len(budgets) == 0 {
			fmt.Println("no budgets configured")
			return
		}
		fmt.Printf("%-20s %-15s %s\n", "Category", "Limit", "Status")
		fmt.Println(repeatStr("-", 60))
		for _, b := range budgets {
			// TODO: calculate actual spending this month
			status := "OK"
			fmt.Printf("%-20s %-15.2f %s\n", b.Category, b.MonthlyLimit, status)
		}
	}
}

func cmdServe(store *storage.JSONStore) {
	srv := api.NewServer(*apiAddr, store)
	fmt.Printf("starting server on %s\n", *apiAddr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `money-backword - household budget CLI

Usage:
  moneyback [flags] <command> [args]

Commands:
  add       Add a transaction (-account, -amount, -category, -desc)
  list      List transactions (-account, -limit)
  accounts  List all accounts
  report    Generate a report (-month, -type)
  budget    Manage budgets (-action set/list/check, -category, -limit)
  serve     Start HTTP API server

Flags:
  -db <path>    Database file (default: ./money.json)
  -api <addr>   API listen address (default: :8080)
  -v            Verbose output

Examples:
  moneyback -db ./ledger.json add -account checking -amount 50.00 -category groceries
  moneyback list -account checking -limit 10
  moneyback report -month 2024-01 -type category
  moneyback budget -action set -category groceries -limit 500

`)
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
