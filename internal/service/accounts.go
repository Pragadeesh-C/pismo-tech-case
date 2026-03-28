package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pragadeesh-c/pismo-tech-case/internal/repository"
)

type CreateAccountInput struct {
	DocumentNumber string
}

type Account struct {
	AccountID      int
	DocumentNumber string
	CreatedAt      time.Time
}

type AccountsService struct {
	repo repository.Querier
}

func NewAccountsService(repo repository.Querier) *AccountsService {
	return &AccountsService{repo: repo}
}

func (s *AccountsService) Create(ctx context.Context, input CreateAccountInput) (*Account, error) {
	if strings.TrimSpace(input.DocumentNumber) == "" {
		return nil, ErrDocNumEmpty
	}

	account, err := s.repo.CreateAccount(ctx, input.DocumentNumber)
	if err != nil {
		if isDocumentConflict(err) {
			return nil, ErrAccountAlreadyExists
		}
		return nil, err
	}

	return &Account{
		AccountID:      int(account.ID),
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt,
	}, nil
}

func isDocumentConflict(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" &&
			pgErr.ConstraintName == "accounts_document_number_key"
	}
	return false
}
