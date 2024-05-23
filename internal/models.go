package internal

import (
	"time"
)

type Account struct {
	AccountId string `json:"accountId"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
}

type Transaction struct {
	TransactionId string    `json:"accountId"`
	ToAccount     string    `json:"receiver"`
	CreatedAt     time.Time `json:"createdAt"`
	Amount        float32   `json:"amount"`
	IsConsumed    bool      `json:"isConsumed"`
}

var Accounts = make(map[string]*Account)
var Transactions = make(map[string]*Transaction)
