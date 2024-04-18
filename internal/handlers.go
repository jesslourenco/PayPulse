package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to GoPay!")
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accs := []*Account{}

	for _, account := range Accounts {
		accs = append(accs, account)
	}

	response := &JsonResponse{Data: &accs}
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusNotFound, Title: "Internal Server Error"}})
		return
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")

	account, found := Accounts[id]
	response := &JsonResponse{Data: &account}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusNotFound, Title: "Account Not Found"}})
		return
	}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusNotFound, Title: "Internal Server Error"}})
	}
}
func PostAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	account := &Account{}

	account.AccountId = params.ByName("account-id")

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request: failure when reading body"}})
		return
	}

	defer r.Body.Close()
	err = json.Unmarshal(body, account)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request: failure when converting to entity"}})
		return
	}

	fmt.Printf("%#v \n", account)

	w.Header().Set("Content-Type", "application/json; charset=UTF8")

	Accounts[account.AccountId] = account
	w.WriteHeader(http.StatusOK)
}
