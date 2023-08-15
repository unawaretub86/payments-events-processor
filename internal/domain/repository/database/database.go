package database

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
)

type (
	Database interface {
		CreatePayment(entities.ProcessPaymentRequest, string) (*entities.ProcessPaymentRequest, error)
		UpdatePayment(string, string) (*string, *string, error)
	}

	databasePayment struct {
		dynamodb dynamodbiface.DynamoDBAPI
		table    string
	}
)

func NewDataBase(dynamodbClient dynamodbiface.DynamoDBAPI) Database {
	const tableName = "payments"

	return &databasePayment{
		dynamodb: dynamodbClient,
		table:    tableName,
	}
}
