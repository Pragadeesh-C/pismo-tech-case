package service

import (
	"context"
	"errors"
	"testing"

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
				Code:           "23505",
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
			mockRepo := &mocks.MockRepo{
				// Inject repository behavior per test case (success or error)
				CreateAccountFn: func(ctx context.Context, documentNumber string) (repository.Account, error) {
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

			// Verify that the service returns the expected error
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("expected %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil && acc == nil {
				t.Fatalf("expected account, got nil")
			}
		})
	}
}
