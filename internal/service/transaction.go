package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pragadeesh-c/pismo-tech-case/internal/constants"
	"github.com/pragadeesh-c/pismo-tech-case/internal/repository"
)

type CreateTransaction struct {
	AccountID     int
	OperationType int
	Amount        float64
}

type Transaction struct {
	ID                int32          `json:"transaction_id" example:"1"`
	AccountID         int32          `json:"account_id" example:"2"`
	OperationTypeName string         `json:"operation_type_name" example:"Purchase with Installments"`
	OperationType     int16          `json:"operation_type" example:"2"`
	Amount            pgtype.Numeric `json:"amount" swaggertype:"number" example:"20.23"`
	EventDate         time.Time      `json:"event_date" example:"2026-03-28T10:00:00Z"`
}

type TransactionService struct {
	repo repository.Querier
}

func NewTransactionService(repo repository.Querier) *TransactionService {
	return &TransactionService{repo: repo}
}

// Create validates the operation type and amount, applies the correct sign convention
// (negative for purchases/withdrawals, positive for credit vouchers), and inserts the transaction into the db.
func (s *TransactionService) Create(ctx context.Context, input CreateTransaction) (*Transaction, error) {
	sign, ok := constants.OperationSign[input.OperationType]
	if !ok {
		return nil, ErrInvalidOperationType
	}

	if input.Amount == 0 {
		return nil, ErrInvalidAmount
	}

	amount := float64(sign) * math.Abs(input.Amount)
	amountInNumeric, err := toNumeric(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount to numeric: %w", err)
	}

	transaction, err := s.repo.CreateTransaction(ctx, repository.CreateTransactionParams{
		AccountID:     int32(input.AccountID),
		OperationType: int16(input.OperationType),
		Amount:        amountInNumeric,
	})

	if err != nil {
		if isReferenceDoesNotExist(err) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return &Transaction{
		ID:                transaction.ID,
		AccountID:         transaction.AccountID,
		OperationTypeName: constants.OperationNames[int(transaction.OperationType)],
		OperationType:     transaction.OperationType,
		Amount:            transaction.Amount,
		EventDate:         transaction.EventDate,
	}, nil
}

// toNumeric converts a float64 to a pgtype.Numeric. This is needed because the repository
// expects a pgtype.Numeric, and the service converts the amount to a float64 before passing it to the repository.
func toNumeric(v float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%f", v))
	return n, err
}

// isReferenceDoesNotExist checks if the error is a foreign key violation.
func isReferenceDoesNotExist(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == PgErrForeignKeyViolation
	}
	return false
}
