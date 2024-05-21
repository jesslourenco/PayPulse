package internal

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("Welcome to GoPay!")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusInternalServerError, Title: "Internal Server Error"}})
		return
	}
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accs := []*Account{}

	for _, account := range Accounts {
		accs = append(accs, account)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	err := json.NewEncoder(w).Encode(&accs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusInternalServerError, Title: "Internal Server Error"}})
		return
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")

	account, found := Accounts[id]
	w.Header().Set("Content-Type", "application/json; charset=UTF8")

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusNotFound, Title: "Account Not Found"}})
		return
	}

	err := json.NewEncoder(w).Encode(&account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusInternalServerError, Title: "Internal Server Error"}})
		return
	}
}
func PostAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")
	_, found := Accounts[id]

	if found {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request: Cannot create new account with this id"}})
		return
	}

	account := &Account{}
	account.AccountId = id

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

	log.Printf("%#v \n", account)

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	Accounts[account.AccountId] = account
}
