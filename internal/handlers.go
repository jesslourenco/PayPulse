package internal

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gopay/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("Welcome to GoPay!")

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
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
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("account-id")

	account, found := Accounts[id]
	w.Header().Set("Content-Type", "application/json; charset=UTF8")

	if !found {
		msg := "Account Not Found"
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg(msg)
		utils.ErrorWithMessage(w, http.StatusNotFound, msg)
		return
	}

	err := json.NewEncoder(w).Encode(&account)
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
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
	err = json.Unmarshal(body, account)

	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Error().Err(err).Msg(err.Error())
		utils.ErrorWithMessage(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	Accounts[account.AccountId] = account
}
