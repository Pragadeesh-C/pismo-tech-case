package service

import (
	"context"
	"errors"
	"fmt"
	"math"
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
		expectedTransaction repository.CreateTransactionRow
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
			expectedTransaction: repository.CreateTransactionRow{
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
		{
			name: "partial discharge",
			input: CreateTransaction{
				AccountID:     1,
				OperationType: 4,
				Amount:        60,
			},
			mockErr:       nil,
			expectedError: nil,
			expectedTransaction: repository.CreateTransactionRow{
				ID:            1,
				AccountID:     1,
				OperationType: 4,
				Amount:        mustNumeric(60),
				Balance:       mustNumeric(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoCalled := false
			mockRepo := &mocks.MockRepo{
				CreateTransactionFn: func(ctx context.Context, transaction repository.CreateTransactionParams) (repository.CreateTransactionRow, error) {
					repoCalled = true

					if tt.mockErr != nil {
						return repository.CreateTransactionRow{}, tt.mockErr
					}

					return repository.CreateTransactionRow{
						ID:            1,
						AccountID:     transaction.AccountID,
						OperationType: transaction.OperationType,
						Amount:        transaction.Amount,
						EventDate:     time.Now(),
					}, nil
				},
				FetchAllDebitTransactionsByAccountIDFn: func(ctx context.Context, accountID int32) ([]repository.Transaction, error) {
					return []repository.Transaction{
						{
							ID:            1,
							AccountID:     1,
							OperationType: 1,
							Amount:        mustNumeric(-50),
							Balance:       mustNumeric(-50),
						},
						{
							ID:            2,
							AccountID:     1,
							OperationType: 1,
							Amount:        mustNumeric(-23.5),
							Balance:       mustNumeric(-23.5),
						},
						{
							ID:            3,
							AccountID:     1,
							OperationType: 1,
							Amount:        mustNumeric(-18.7),
							Balance:       mustNumeric(-18.7),
						},
					}, nil
				},
				UpdateTransactionByIDFn: func(ctx context.Context, arg repository.UpdateTransactionByIDParams) error {
					return nil
				},
			}

			mockStore := &mocks.MockStore{
				WithTransactionFn: func(ctx context.Context, fn func(q repository.Querier) error) error {
					return fn(mockRepo)
				},
			}

			service := NewTransactionService(mockStore)

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

			if tt.expectedError != nil {
				return
			}

			txAmountInNumeric, err := numericToFloat(transaction.Amount)
			if err != nil {
				t.Fatalf("error occured while converting transaction amount to numeric: %v", err)
			}
			txAmountAbs := math.Abs(*txAmountInNumeric)

			expectedAmountInNumeric, err := numericToFloat(tt.expectedTransaction.Amount)
			if err != nil {
				t.Fatalf("error occured while converting expected transaction amount to numeric: %v", err)
			}
			expectedAmountAbs := math.Abs(*expectedAmountInNumeric)

			if txAmountAbs != expectedAmountAbs {
				t.Fatalf("expected transaction amount: %.2f, got: %.2f", expectedAmountAbs, txAmountAbs)
			}

			txBalanceInNumeric, err := numericToFloat(transaction.Balance)
			if err != nil {
				t.Fatalf("error occured while converting transaction balance to numeric: %v", err)
			}
			txBalanceAbs := math.Abs(*txBalanceInNumeric)

			expectedBalanceInNumeric, err := numericToFloat(tt.expectedTransaction.Balance)
			if err != nil {
				t.Fatalf("error occured while converting expected transaction balance to numeric: %v", err)
			}
			expectedBalanceAbs := math.Abs(*expectedBalanceInNumeric)

			if txBalanceAbs != expectedBalanceAbs {
				t.Fatalf("expected transaction balance: %.2f, got: %.2f", expectedBalanceAbs, txBalanceAbs)
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
