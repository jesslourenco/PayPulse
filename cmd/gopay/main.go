package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gopay/internal"
)

func initDB() {
	// temp fake db for accounts and transactions
	internal.Accounts["0001"] = &internal.Account{
		AccountId: "0001",
		Name:      "Shankar",
		LastName:  "Nakai",
	}

	internal.Accounts["0002"] = &internal.Account{
		AccountId: "0002",
		Name:      "Jessica",
		LastName:  "Lourenco",
	}

	internal.Accounts["0003"] = &internal.Account{
		AccountId: "0003",
		Name:      "Caio",
		LastName:  "Henrique",
	}

	internal.Accounts["0004"] = &internal.Account{
		AccountId: "0004",
		Name:      "Karina",
		LastName:  "Domingues",
	}

	internal.Transactions["1000000"] = &internal.Transaction{
		TransactionId: "1000000",
		Owner:         "0001",
		Sender:        "0001",
		Receiver:      "0001",
		CreatedAt:     time.Now(),
		Amount:        7000.00,
		IsConsumed:    false,
	}

	internal.Transactions["2000000"] = &internal.Transaction{
		TransactionId: "2000000",
		Owner:         "0002",
		Sender:        "0002",
		Receiver:      "0002",
		CreatedAt:     time.Now(),
		Amount:        3000.00,
		IsConsumed:    false,
	}
}

func main() {
	router := internal.Router(internal.Routes())
	initDB()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("Server started at port :8080")

	log.
		Fatal().
		Err(http.ListenAndServe(":8080", router)).
		Msg("Server closed")
}
