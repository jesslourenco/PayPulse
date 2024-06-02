package internal

import "github.com/julienschmidt/httprouter"

type Route struct {
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

func Routes() []Route {
	return []Route{
		{"GET", "/", Index},
		{"GET", "/accounts", GetAllAccounts},
		{"GET", "/accounts/:account-id", GetAccount},
		{"POST", "/accounts/:account-id", PostAccount},
		{"GET", "/accounts/:account-id/transactions", GetAllTransactions},
		{"GET", "/transactions/:transaction-id", GetTransaction},
		{"POST", "/transactions/:transaction-id", PostTransaction},
		{"GET", "/accounts/:account-id/balance", GetBalance},
	}
}
