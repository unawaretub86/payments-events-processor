package repository

import (
	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
	"github.com/unawaretub86/payments-events-processor/internal/domain/repository/database"
)

type (
	RepositoryPayment interface {
		CreatePayment(entities.ProcessPaymentRequest, string) (*entities.ProcessPaymentRequest, error)
		UpdatePayment(string, string) (*string, *string, error)
	}

	repositoryPayment struct {
		database database.Database
	}
)

func NewRepository(database database.Database) RepositoryPayment {
	return &repositoryPayment{
		database: database,
	}
}
