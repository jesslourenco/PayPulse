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
		{"POST", "/account/:account-id", PostAccount},
	}
}
