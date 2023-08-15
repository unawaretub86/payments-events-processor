package repository

import "github.com/unawaretub86/payments-events-processor/internal/domain/entities"

func (repositoryPayment repositoryPayment) CreatePayment(order entities.ProcessPaymentRequest, requestId string) (*entities.ProcessPaymentRequest, error) {
	return repositoryPayment.database.CreatePayment(order, requestId)
}

func (repositoryPayment repositoryPayment) UpdatePayment(orderId, requestId string) (*string, *string, error) {
	return repositoryPayment.database.UpdatePayment(orderId, requestId)
}
