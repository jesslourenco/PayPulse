package utils

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gopay/internal/models"
	"github.com/rs/zerolog/log"
)

var (
	ErrMaxAttemps = errors.New("operation failed after maximum attempts")
	ErrNotAFunc   = errors.New("fn must be a function")
)

func ErrorWithMessage(w http.ResponseWriter, status int, message string) {
	resp := ErrorResponse{
		Status:  status,
		Message: message,
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		WithPayload(w, http.StatusUnprocessableEntity, []byte(`{"error": "Unable to Return Payload"}`))
	}
	WithPayload(w, status, payload)
}

func WithPayload(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

func Retry(fn func() error, op string) error {
	const maxAttempts = 5
	const delay = 2 * time.Second

	log.Info().Msgf("Operation %v failed. Retrying...", op)

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err := fn()
		if err == nil {
			return nil
		}
		log.Info().Msgf("Attempt %d failed: %v. Retrying...\n", attempts, err)
		time.Sleep(delay)
	}

	return ErrMaxAttemps
}

func Deposit(transaction *models.Transaction) bool {
	models.Transactions[transaction.TransactionId] = transaction
	return true
}

func Withdrawal(transaction *models.Transaction) bool {
	if !Debit(transaction) {
		return false
	}

	transaction.IsConsumed = true
	models.Transactions[transaction.TransactionId] = transaction

	return true
}

func Pay(transaction *models.Transaction) bool {
	if !Debit(transaction) {
		return false
	}

	transaction.IsConsumed = true
	models.Transactions[transaction.TransactionId] = transaction

	Credit(transaction)

	return true
}

func Debit(transaction *models.Transaction) bool {
	var balance float64
	var oldest *models.Transaction

	for _, t := range models.Transactions {
		if t.Owner == transaction.Owner && !t.IsConsumed {
			balance += float64(t.Amount)
			if oldest == nil || t.CreatedAt.Before(oldest.CreatedAt) {
				oldest = t
			}
		}
	}

	if (balance + float64(transaction.Amount)) < 0 {
		return false
	}

	oldest.IsConsumed = true

	if (oldest.Amount + transaction.Amount) != 0 {
		id := GetTransactionUUID()

		new := &models.Transaction{
			TransactionId: id,
			Owner:         transaction.Owner,
			Sender:        transaction.Owner,
			Receiver:      transaction.Owner,
			CreatedAt:     time.Now(),
			Amount:        oldest.Amount + transaction.Amount,
			IsConsumed:    false,
		}
		models.Transactions[new.TransactionId] = new
	}

	return true
}

func Credit(transaction *models.Transaction) {
	id := GetTransactionUUID()
	t := &models.Transaction{
		TransactionId: id,
		Owner:         transaction.Receiver,
		Sender:        transaction.Owner,
		Receiver:      transaction.Receiver,
		CreatedAt:     time.Now(),
		Amount:        float32(math.Abs(float64(transaction.Amount))),
		IsConsumed:    false,
	}
	models.Transactions[t.TransactionId] = t
}

func GetAccountUUID() string {
	id := uuid.NewString()

	_, found := models.Accounts[id]
	for found {
		id = uuid.NewString()
		_, found = models.Accounts[id]
	}

	return id
}

func GetTransactionUUID() string {
	id := uuid.NewString()

	_, found := models.Transactions[id]
	for found {
		id = uuid.NewString()
		_, found = models.Transactions[id]
	}

	return id
}
