package mocks

import (
	"context"

	"github.com/pragadeesh-c/pismo-tech-case/internal/repository"
)

// MockRepo is a test implementation of repository.Querier.
// It enables injecting custom behavior for repository calls,
// allowing the service layer to be tested in isolation.
type MockRepo struct {
	CreateAccountFn func(ctx context.Context, documentNumber string) (repository.Account, error)
	GetAccountFn    func(ctx context.Context, accountId int32) (repository.Account, error)
}

func (m *MockRepo) CreateAccount(ctx context.Context, documentNumber string) (repository.Account, error) {
	return m.CreateAccountFn(ctx, documentNumber)
}

func (m *MockRepo) GetAccount(ctx context.Context, accountId int32) (repository.Account, error) {
	return m.GetAccountFn(ctx, accountId)
}
