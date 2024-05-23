package internal

import (
	"io"
	"net/http"

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
