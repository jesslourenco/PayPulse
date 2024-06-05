package internal

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/utils"
	jsoniter "github.com/json-iterator/go"

	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

const AccountIdParam = "account-id"
const TransactionIdParam = "transaction-id"
const OneMegabyte = 1048576

var ErrTransactionNotFound = errors.New("transaction not found")
var ErrAccountNotFound = errors.New("account not found")
var ErrReceiverNotFound = errors.New("receiver account not found")
var ErrSenderNotFound = errors.New("sender account not found")
var ErrInsufficentBalance = errors.New("insufficient balance")

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF8")

	res, err := jsoniter.Marshal("Welcome to GoPay!")

	if err != nil {
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
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WithPayload(w, http.StatusOK, res)
}

func GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName(AccountIdParam)

	account, found := models.Accounts[id]

	if !found {
		log.Error().Err(ErrAccountNotFound).Msg("Handler::GetAccount")
		utils.ErrorWithMessage(w, http.StatusNotFound, ErrAccountNotFound.Error())
		return
	}

	res, err := jsoniter.Marshal(&account)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WithPayload(w, http.StatusOK, res)
}

func PostAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName(AccountIdParam)
	_, found := models.Accounts[id]

	if found {
		msg := "Cannot create new account with this id (duplicate)"
		log.Error().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusBadRequest, msg)
		return
	}

	account := &models.Account{}
	account.AccountId = id

	body, err := io.ReadAll(io.LimitReader(r.Body, OneMegabyte))
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()
	err = jsoniter.Unmarshal(body, &account)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	models.Accounts[account.AccountId] = account
	utils.WithPayload(w, http.StatusCreated, nil)
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	transactions := []*models.Transaction{}

	accountId := params.ByName(AccountIdParam)

	_, found := models.Accounts[accountId]

	if !found {
		log.Error().Err(ErrAccountNotFound).Msg("Handler::GetAllTransactions")
		utils.ErrorWithMessage(w, http.StatusNotFound, ErrAccountNotFound.Error())
		return
	}

	for _, transaction := range models.Transactions {
		if transaction.Owner == accountId {
			transactions = append(transactions, transaction)
		}
	}

	res, err := jsoniter.Marshal(&transactions)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WithPayload(w, http.StatusOK, res)
}

func GetTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName(TransactionIdParam)

	transaction, found := models.Transactions[id]

	if !found {
		log.Error().Err(ErrTransactionNotFound).Msg("Handler::GetTransaction")
		utils.ErrorWithMessage(w, http.StatusNotFound, ErrTransactionNotFound.Error())
		return
	}

	res, err := jsoniter.Marshal(&transaction)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WithPayload(w, http.StatusOK, res)
}

func PostTransaction(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName(TransactionIdParam)
	_, found := models.Transactions[id]

	if found {
		msg := "Cannot create new transaction with this id (duplicate)"
		log.Error().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusBadRequest, msg)
		return
	}

	transaction := &models.Transaction{}
	transaction.TransactionId = id
	transaction.CreatedAt = time.Now()
	transaction.IsConsumed = false

	body, err := io.ReadAll(io.LimitReader(r.Body, OneMegabyte))
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()
	err = jsoniter.Unmarshal(body, &transaction)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, found = models.Accounts[transaction.Receiver]
	if !found {
		log.Error().Err(ErrReceiverNotFound).Msg("Handler::PostTransaction")
		utils.ErrorWithMessage(w, http.StatusNotFound, ErrReceiverNotFound.Error())
		return
	}

	_, found = models.Accounts[transaction.Receiver]
	if !found {
		log.Error().Err(ErrSenderNotFound).Msg("Handler::PostTransaction")
		utils.ErrorWithMessage(w, http.StatusNotFound, ErrSenderNotFound.Error())
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
			log.Error().Err(ErrInsufficentBalance).Msg("Handler::PostTransaction")
			utils.ErrorWithMessage(w, http.StatusForbidden, ErrInsufficentBalance.Error())
			return
		}

		utils.WithPayload(w, http.StatusCreated, nil)
		return
	}

	success = utils.Pay(transaction)
	if !success {
		log.Error().Err(ErrInsufficentBalance).Msg("Handler::PostTransaction")
		utils.ErrorWithMessage(w, http.StatusForbidden, ErrInsufficentBalance.Error())
		return
	}

	utils.WithPayload(w, http.StatusCreated, nil)
}

func GetBalance(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName(AccountIdParam)

	_, found := models.Accounts[id]

	if !found {
		log.Error().Err(ErrAccountNotFound).Msg("Handler::GetBalance")
		utils.ErrorWithMessage(w, http.StatusNotFound, ErrAccountNotFound.Error())
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
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WithPayload(w, http.StatusOK, res)
}
