package database

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type (
	Database interface {
		CreatePayment(string, string) (*string, error)
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
