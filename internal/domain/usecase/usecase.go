package usecase

import (
	"os"

	"github.com/unawaretub86/payments-events-processor/internal/domain/repository"
)

type (
	UseCase interface {
		CreatePayment(string, string) (*string, error)
	}

	useCase struct {
		repositoryPayment repository.RepositoryPayment
		queueURL          string
	}
)

func NewUseOrder(repositoryPayment repository.RepositoryPayment) UseCase {
	queueURL := os.Getenv("SQS_URL")

	return &useCase{
		repositoryPayment: repositoryPayment,
		queueURL:          queueURL,
	}
}
