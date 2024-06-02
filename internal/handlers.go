package internal

import (
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/gopay/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF8")

	res, err := jsoniter.Marshal("Welcome to GoPay!")

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WithPayload(w, http.StatusOK, res)
}

func GetAllAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	accs := []*Account{}

	for _, account := range Accounts {
		accs = append(accs, account)
	}

	res, err := jsoniter.Marshal(&accs)

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WithPayload(w, http.StatusOK, res)
}

func GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")

	account, found := Accounts[id]

	if !found {
		msg := "Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	res, err := jsoniter.Marshal(&account)
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WithPayload(w, http.StatusOK, res)
}

func PostAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")
	_, found := Accounts[id]

	if found {
		msg := "Cannot create new account with this id (duplicate)"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusBadRequest, msg)
		return
	}

	account := &Account{}
	account.AccountId = id

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()
	err = jsoniter.Unmarshal(body, &account)

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	Accounts[account.AccountId] = account
	utils.WithPayload(w, http.StatusCreated, nil)
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	transactions := []*Transaction{}

	accountId := params.ByName("account-id")

	_, found := Accounts[accountId]

	if !found {
		msg := "Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	for _, transaction := range Transactions {
		if transaction.Sender == accountId {
			transactions = append(transactions, transaction)
		}
	}

	res, err := jsoniter.Marshal(&transactions)

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WithPayload(w, http.StatusOK, res)
}

func GetTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("transaction-id")

	transaction, found := Transactions[id]

	if !found {
		msg := "Transaction Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	res, err := jsoniter.Marshal(&transaction)
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WithPayload(w, http.StatusOK, res)
}

func PostTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("transaction-id")
	_, found := Transactions[id]

	if found {
		msg := "Cannot create new transaction with this id (duplicate)"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusBadRequest, msg)
		return
	}

	transaction := &Transaction{}
	transaction.TransactionId = id
	transaction.CreatedAt = time.Now()
	transaction.IsConsumed = false

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()
	err = jsoniter.Unmarshal(body, &transaction)

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, found = Accounts[transaction.Receiver]
	if !found {
		msg := "Receiver Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	_, found = Accounts[transaction.Receiver]
	if !found {
		msg := "Sender Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	Transactions[transaction.TransactionId] = transaction
	utils.WithPayload(w, http.StatusCreated, nil)

	// TODO: Add account balance validation
}

func GetBalance(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")

	_, found := Accounts[id]

	if !found {
		msg := "Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	balance := Balance{
		AccountId: id,
		Amount:    0.00,
	}

	for _, transaction := range Transactions {
		if transaction.Owner == id && !transaction.IsConsumed {
			balance.Amount += float64(transaction.Amount)
		}
	}

	res, err := jsoniter.Marshal(&balance)
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WithPayload(w, http.StatusOK, res)
}
