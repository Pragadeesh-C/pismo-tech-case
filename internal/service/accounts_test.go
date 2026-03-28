package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pragadeesh-c/pismo-tech-case/internal/mocks"
	"github.com/pragadeesh-c/pismo-tech-case/internal/repository"
)

// TestCreateAccount verifies the AccountsService.Create method.
// It covers validation, duplicate handling (unique constraint),
// successful creation, and unexpected repository errors.
func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateAccountInput
		mockErr       error
		expectedError error
	}{
		{
			name: "success",
			input: CreateAccountInput{
				DocumentNumber: "123",
			},
			mockErr:       nil,
			expectedError: nil,
		},
		{
			name: "document number has only empty spaces",
			input: CreateAccountInput{
				DocumentNumber: " ",
			},
			mockErr:       nil,
			expectedError: ErrDocNumEmpty,
		},
		{
			name: "document number empty",
			input: CreateAccountInput{
				DocumentNumber: "",
			},
			mockErr: nil,
			// Service trims input and should treat this as empty
			expectedError: ErrDocNumEmpty,
		},
		{
			name: "duplicate account",
			input: CreateAccountInput{
				DocumentNumber: "123",
			},
			// Simulates Postgres unique constraint violation (23505)
			mockErr: &pgconn.PgError{
				Code:           PgErrUniqueViolation,
				ConstraintName: "accounts_document_number_key",
			},
			expectedError: ErrAccountAlreadyExists,
		},
		{
			name: "unexpected db error",
			input: CreateAccountInput{
				DocumentNumber: "123",
			},
			mockErr:       ErrDBDown,
			expectedError: ErrDBDown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoCalled := false
			mockRepo := &mocks.MockRepo{
				// Inject repository behavior per test case (success or error)
				CreateAccountFn: func(ctx context.Context, documentNumber string) (repository.Account, error) {
					repoCalled = true
					if tt.mockErr != nil {
						return repository.Account{}, tt.mockErr
					}
					return repository.Account{
						ID:             1,
						DocumentNumber: documentNumber,
					}, nil
				},
			}

			service := NewAccountsService(mockRepo)

			acc, err := service.Create(context.Background(), tt.input)

			if strings.TrimSpace(tt.input.DocumentNumber) == "" && repoCalled {
				t.Fatal("repo called for invalid input")
			}

			if strings.TrimSpace(tt.input.DocumentNumber) != "" && tt.expectedError == nil && !repoCalled {
				t.Fatal("repo should be called for valid input")
			}

			// Verify that the service returns the expected error
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("expected %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil {
				if acc == nil {
					t.Fatalf("expected account, got nil")
				}

				if acc.DocumentNumber != strings.TrimSpace(tt.input.DocumentNumber) {
					t.Fatalf("expected document number %v, got %v",
						strings.TrimSpace(tt.input.DocumentNumber),
						acc.DocumentNumber)
				}
			}
		})
	}
}

// TestGetAccount verifies the AccountsService.GetAccountByID method.
// It covers ID validation, not-found handling,
// successful retrieval, and unexpected repository errors.
func TestGetAccount(t *testing.T) {
	tests := []struct {
		name            string
		accountID       int
		mockErr         error
		expectedAccount *Account
		expectedError   error
	}{
		{
			name:          "success",
			accountID:     16,
			mockErr:       nil,
			expectedError: nil,
			expectedAccount: &Account{
				AccountID:      16,
				DocumentNumber: "123",
			},
		},
		{
			name:            "account id invalid",
			accountID:       -1,
			mockErr:         nil,
			expectedError:   ErrInvalidAccountID,
			expectedAccount: nil,
		},
		{
			name:            "zero account id",
			accountID:       0,
			mockErr:         nil,
			expectedError:   ErrInvalidAccountID,
			expectedAccount: nil,
		},
		{
			name:            "account not found",
			accountID:       2,
			mockErr:         pgx.ErrNoRows,
			expectedError:   ErrAccountNotFound,
			expectedAccount: nil,
		},
		{
			name:            "unexpected db error",
			accountID:       5,
			mockErr:         ErrDBDown,
			expectedError:   ErrDBDown,
			expectedAccount: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoCalled := false
			mockRepo := &mocks.MockRepo{
				GetAccountFn: func(ctx context.Context, accountId int32) (repository.Account, error) {
					repoCalled = true
					if tt.mockErr != nil {
						return repository.Account{}, tt.mockErr
					}
					return repository.Account{
						ID:             accountId,
						DocumentNumber: "123",
						CreatedAt:      time.Now(),
					}, nil
				},
			}

			service := NewAccountsService(mockRepo)

			account, err := service.GetAccountByID(context.Background(), tt.accountID)

			if repoCalled && tt.accountID < 1 {
				t.Fatal("repo should not be called for invalid input")
			}

			if tt.accountID > 0 && tt.expectedError == nil && !repoCalled {
				t.Fatal("repo should be called for valid input")
			}

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("expected %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil && account == nil {
				t.Fatalf("expected account, got nil")
			}

			if tt.expectedError == nil {
				if account.AccountID != tt.expectedAccount.AccountID || account.DocumentNumber != tt.expectedAccount.DocumentNumber {
					t.Fatalf("expected %v got %v", tt.expectedAccount, account)
				}
				if account.CreatedAt.IsZero() {
					t.Fatal("expected CreatedAt to be set")
				}
			}

		})
	}
}
