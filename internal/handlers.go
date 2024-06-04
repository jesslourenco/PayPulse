package internal

import (
	"io"
	"net/http"
	"time"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/utils"
	jsoniter "github.com/json-iterator/go"

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
	accs := []*models.Account{}

	for _, account := range models.Accounts {
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

	account, found := models.Accounts[id]

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
	_, found := models.Accounts[id]

	if found {
		msg := "Cannot create new account with this id (duplicate)"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusBadRequest, msg)
		return
	}

	account := &models.Account{}
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

	models.Accounts[account.AccountId] = account
	utils.WithPayload(w, http.StatusCreated, nil)
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	transactions := []*models.Transaction{}

	accountId := params.ByName("account-id")

	_, found := models.Accounts[accountId]

	if !found {
		msg := "Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	for _, transaction := range models.Transactions {
		if transaction.Owner == accountId {
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

	transaction, found := models.Transactions[id]

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
	_, found := models.Transactions[id]

	if found {
		msg := "Cannot create new transaction with this id (duplicate)"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusBadRequest, msg)
		return
	}

	transaction := &models.Transaction{}
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

	_, found = models.Accounts[transaction.Receiver]
	if !found {
		msg := "Receiver Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	_, found = models.Accounts[transaction.Receiver]
	if !found {
		msg := "Sender Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	transaction.Owner = transaction.Sender
	var success bool

	if transaction.Sender == transaction.Receiver {
		if transaction.Amount > 0 {
			success = utils.Deposit(transaction)
		} else {
			success = utils.Withdrawal(transaction)
		}

		if !success {
			msg := "Insufficient Balance"
			zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
			log.Info().Msg(msg)
			utils.ErrorWithMessage(w, http.StatusForbidden, msg)
			return
		}

		utils.WithPayload(w, http.StatusCreated, nil)
		return
	}

	success = utils.Pay(transaction)
	if !success {
		msg := "Insufficient Balance"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusForbidden, msg)
		return
	}

	utils.WithPayload(w, http.StatusCreated, nil)
}

func GetBalance(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")

	_, found := models.Accounts[id]

	if !found {
		msg := "Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	balance := models.Balance{
		AccountId: id,
		Amount:    0.00,
	}

	for _, transaction := range models.Transactions {
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
