package internal

import (
	"time"
)

type Account struct {
	AccountId string `json:"account-id"`
	Name      string `json:"name"`
	LastName  string `json:"lastname"`
}

type Transaction struct {
	TransactionId string    `json:"account-id"`
	ToAccount     string    `json:"receiver"`
	CreatedAt     time.Time `json:"createdAt"`
	Amount        float32   `json:"amount"`
	IsConsumed    bool      `json:"isConsumed"`
}

var Accounts = make(map[string]*Account)
var Transactions = make(map[string]*Transaction)
