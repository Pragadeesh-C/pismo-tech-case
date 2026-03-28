package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pragadeesh-c/pismo-tech-case/internal/mocks"
	"github.com/pragadeesh-c/pismo-tech-case/internal/repository"
)

// TestCreateTransaction verifies the TransactionService.Create method.
// It covers operation type validation, amount validation,
// foreign key (account not found) handling, and successful creation.
func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                string
		input               CreateTransaction
		mockErr             error
		expectedError       error
		expectedTransaction repository.CreateTransactionParams
	}{
		{
			name: "success",
			input: CreateTransaction{
				AccountID:     1,
				OperationType: 1,
				Amount:        20.3,
			},
			mockErr:       nil,
			expectedError: nil,
			expectedTransaction: repository.CreateTransactionParams{
				AccountID:     1,
				OperationType: 1,
				Amount:        mustNumeric(20.3),
			},
		},
		{
			name: "account not found",
			input: CreateTransaction{
				AccountID:     20,
				OperationType: 2,
				Amount:        10.23,
			},
			mockErr: &pgconn.PgError{
				Code: PgErrForeignKeyViolation,
			},
			expectedError: ErrAccountNotFound,
		},
		{
			name: "invalid operation type",
			input: CreateTransaction{
				AccountID:     1,
				OperationType: 0,
				Amount:        -24.30,
			},
			mockErr:       nil,
			expectedError: ErrInvalidOperationType,
		},
		{
			name: "invalid amount(0)",
			input: CreateTransaction{
				AccountID:     1,
				OperationType: 2,
				Amount:        0,
			},
			mockErr:       nil,
			expectedError: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoCalled := false
			mockRepo := &mocks.MockRepo{
				CreateTransactionFn: func(ctx context.Context, transaction repository.CreateTransactionParams) (repository.Transaction, error) {
					repoCalled = true

					if tt.mockErr != nil {
						return repository.Transaction{}, tt.mockErr
					}

					return repository.Transaction{
						ID:            1,
						AccountID:     transaction.AccountID,
						OperationType: transaction.OperationType,
						Amount:        transaction.Amount,
						EventDate:     time.Now(),
					}, nil
				},
			}

			service := NewTransactionService(mockRepo)

			transaction, err := service.Create(context.Background(), tt.input)

			if repoCalled && tt.input.Amount == 0 {
				t.Fatalf("repo should not be called for invalid input %v %s", tt.input.Amount, tt.name)
			}

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("expected %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil && transaction == nil {
				t.Fatal("expected transaction, got nil")
			}
		})
	}
}

// mustNumeric converts a float64 to a pgtype.Numeric. This is needed because the repository
// expects a pgtype.Numeric, and the service converts the amount to a float64 before passing it to the repository.
func mustNumeric(v float64) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(fmt.Sprintf("%f", v)); err != nil {
		panic(err)
	}
	return n
}
