package main

import (
	"fmt"
	"errors"
	"time"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Account interface defines the basic operations for bank accounts
type Account interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	CheckBalance() float64
	AccountType() string
	GetTransactionHistory() []string
}

// Transaction represents a banking transaction
type Transaction struct {
	Type      string
	Amount    float64
	Timestamp time.Time
	Status    string
}

// BaseAccount contains common fields for all account types
type BaseAccount struct {
	balance            *float64
	transactionHistory []string
}

// SavingsAccount implements the Account interface
type SavingsAccount struct {
	BaseAccount
	withdrawalLimit float64
}

// CurrentAccount implements the Account interface with overdraft facility
type CurrentAccount struct {
	BaseAccount
	overdraftLimit float64
}

// NewSavingsAccount creates a new savings account with initial balance
func NewSavingsAccount(initialBalance float64) *SavingsAccount {
	balance := initialBalance
	return &SavingsAccount{
		BaseAccount: BaseAccount{
			balance:            &balance,
			transactionHistory: make([]string, 0),
		},
		withdrawalLimit: 500.0,
	}
}

// NewCurrentAccount creates a new current account with initial balance
func NewCurrentAccount(initialBalance float64) *CurrentAccount {
	balance := initialBalance
	return &CurrentAccount{
		BaseAccount: BaseAccount{
			balance:            &balance,
			transactionHistory: make([]string, 0),
		},
		overdraftLimit: 200.0,
	}
}

// Deposit adds money to the account
func (a *BaseAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	*a.balance += amount
	a.recordTransaction("Deposit", amount, "Success")
	return nil
}

// Withdraw removes money from the savings account
func (s *SavingsAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > s.withdrawalLimit {
		return fmt.Errorf("withdrawal amount exceeds limit of ₹%.2f", s.withdrawalLimit)
	}
	if amount > *s.balance {
		return errors.New("insufficient funds")
	}
	*s.balance -= amount
	s.recordTransaction("Withdrawal", amount, "Success")
	return nil
}

// Withdraw removes money from the current account
func (c *CurrentAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > *c.balance+c.overdraftLimit {
		return errors.New("amount exceeds available balance and overdraft limit")
	}
	*c.balance -= amount
	c.recordTransaction("Withdrawal", amount, "Success")
	return nil
}

// CheckBalance returns the current balance
func (a *BaseAccount) CheckBalance() float64 {
	return *a.balance
}

// AccountType returns "Savings" for savings account
func (s *SavingsAccount) AccountType() string {
	return "Savings"
}

// AccountType returns "Current" for current account
func (c *CurrentAccount) AccountType() string {
	return "Current"
}

// recordTransaction adds a transaction to the history
func (a *BaseAccount) recordTransaction(transType string, amount float64, status string) {
	transaction := fmt.Sprintf("%s - Type: %s, Amount: ₹%.2f, Balance: ₹%.2f",
		time.Now().Format("2006-01-02 15:04:05"),
		transType,
		amount,
		*a.balance)
	a.transactionHistory = append(a.transactionHistory, transaction)
}

// GetTransactionHistory returns the transaction history
func (a *BaseAccount) GetTransactionHistory() []string {
	return a.transactionHistory
}

// Global map to store accounts
var accounts = make(map[string]Account)

// AccountRequest represents the JSON request for creating an account
type AccountRequest struct {
	AccountType    int     `json:"accountType"`
	InitialBalance float64 `json:"initialBalance"`
	AccountId      string  `json:"accountId"`
}

// AmountRequest represents the JSON request for deposit/withdrawal
type AmountRequest struct {
	Amount float64 `json:"amount"`
}

// AccountResponse represents the JSON response for account operations
type AccountResponse struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

// createAccountHandler handles the creation of new accounts
func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, exists := accounts[req.AccountId]; exists {
		http.Error(w, "Account ID already exists", http.StatusBadRequest)
		return
	}

	var account Account
	switch req.AccountType {
	case 1:
		account = NewSavingsAccount(req.InitialBalance)
	case 2:
		account = NewCurrentAccount(req.InitialBalance)
	default:
		http.Error(w, "Invalid account type", http.StatusBadRequest)
		return
	}

	accounts[req.AccountId] = account

	response := AccountResponse{
		ID:      req.AccountId,
		Type:    account.AccountType(),
		Balance: account.CheckBalance(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getAccountsHandler returns a list of all accounts
func getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var accountList []AccountResponse
	for id, account := range accounts {
		accountList = append(accountList, AccountResponse{
			ID:      id,
			Type:    account.AccountType(),
			Balance: account.CheckBalance(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountList)
}

// depositHandler handles deposits to an account
func depositHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]

	account, exists := accounts[accountId]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var req AmountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := account.Deposit(req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := AccountResponse{
		ID:      accountId,
		Type:    account.AccountType(),
		Balance: account.CheckBalance(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// withdrawHandler handles withdrawals from an account
func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]

	account, exists := accounts[accountId]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var req AmountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := account.Withdraw(req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := AccountResponse{
		ID:      accountId,
		Type:    account.AccountType(),
		Balance: account.CheckBalance(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getTransactionHistoryHandler returns the transaction history for an account
func getTransactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["id"]

	account, exists := accounts[accountId]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account.GetTransactionHistory())
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// API routes
	r.HandleFunc("/api/accounts", createAccountHandler).Methods("POST")
	r.HandleFunc("/api/accounts", getAccountsHandler).Methods("GET")
	r.HandleFunc("/api/accounts/{id}/deposit", depositHandler).Methods("POST")
	r.HandleFunc("/api/accounts/{id}/withdraw", withdrawHandler).Methods("POST")
	r.HandleFunc("/api/accounts/{id}/history", getTransactionHistoryHandler).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", c.Handler(r))
}