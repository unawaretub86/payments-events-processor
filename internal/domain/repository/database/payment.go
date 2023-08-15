package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (d *databasePayment) CreatePayment(orderId, requestId string) (*string, error) {

	status := "PENDING"

	item := map[string]*dynamodb.AttributeValue{
		"orderId": {S: aws.String(orderId)},
		"Status":  {S: aws.String(status)},
	}

	_, err := d.dynamodb.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      item,
	})

	fmt.Printf("[RequestId: %s], [PutItem result: %v]", orderId, requestId)

	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", err, requestId)
		return nil, err
	}

	return &orderId, nil
}
