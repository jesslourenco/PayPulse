package internal

import "github.com/julienschmidt/httprouter"

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

func Routes() []Route {
	return []Route{
		{"Index", "GET", "/", Index},
		{"AccountIndex", "GET", "/accounts", GetAllAccounts},
		{"Account", "GET", "/accounts/:account-id", GetAccount},
		{"AccountPost", "POST", "/account/:account-id", PostAccount},
	}
}
