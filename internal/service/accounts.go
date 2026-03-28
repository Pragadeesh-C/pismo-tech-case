package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pragadeesh-c/pismo-tech-case/internal/repository"
)

type CreateAccountInput struct {
	DocumentNumber string
}

type Account struct {
	AccountID      int       `json:"accountID" example:"123"`
	DocumentNumber string    `json:"document_number" example:"12345"`
	CreatedAt      time.Time `json:"createdAt" example:"2026-03-28T10:00:00Z"`
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

func (s *AccountsService) GetAccountByID(ctx context.Context, id int) (*Account, error) {
	if id <= 0 {
		return nil, ErrInvalidAccountID
	}

	account, err := s.repo.GetAccount(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
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
