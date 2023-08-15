package usecase

import (
	"os"

	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
	"github.com/unawaretub86/payments-events-processor/internal/domain/repository"
)

type (
	UseCase interface {
		CreatePayment(string, string) (*entities.ProcessPaymentRequest, error)
		UpdatePayment(string, string) error
	}

	useCase struct {
		repositoryPayment repository.RepositoryPayment
		queueURL          string
	}
)

func NewUsePayment(repositoryPayment repository.RepositoryPayment) UseCase {
	queueURL := os.Getenv("SQS_URL")

	return &useCase{
		repositoryPayment: repositoryPayment,
		queueURL:          queueURL,
	}
}
