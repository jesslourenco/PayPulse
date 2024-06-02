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
	TransactionId string    `json:"transactionId"`
	Owner         string    `json:"owner"`
	Sender        string    `json:"sender"`
	Receiver      string    `json:"receiver"`
	CreatedAt     time.Time `json:"createdAt"`
	Amount        float32   `json:"amount"`
	IsConsumed    bool      `json:"isConsumed"`
}

type Balance struct {
	AccountId string  `json:"accountId"`
	Amount    float64 `json:"balance"`
}

var Accounts = make(map[string]*Account)
var Transactions = make(map[string]*Transaction)
