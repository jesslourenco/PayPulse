package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var accounts = make(map[string]*Account)

// var transactions = make(map[string]*internal.Transaction)

type JsonResponse struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int16  `json:"status"`
	Title  string `json:"title"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to GoPay!")
}

// GET /accounts
func AccIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accs := []*Account{}

	accounts["0001"] = &Account{
		AccountId: "0001",
		Name:      "Shankar",
		LastName:  "Nakai",
	}

	for _, account := range accounts {
		accs = append(accs, account)
	}

	response := &JsonResponse{Data: &accs}
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
