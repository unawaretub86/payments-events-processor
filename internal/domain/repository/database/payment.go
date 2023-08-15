package database

import (
	"fmt"
	"log"

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

	fmt.Printf("[RequestId: %s], [PutItem result: %v]", requestId, orderId)

	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
		return nil, err
	}

	return &orderId, nil
}

func (d *databasePayment) UpdatePayment(orderId, requestId string) (*string, *string, error) {

	paid := "PAID"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#Y": aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":y": {
				S: aws.String(paid),
			},
		},
		TableName: aws.String(d.table),
		Key: map[string]*dynamodb.AttributeValue{
			"orderId": {
				S: aws.String(orderId),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET #Y = :y"),
	}

	_, err := d.dynamodb.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	fmt.Printf("[RequestId: %s], [UpdateItem result: %v]", requestId, orderId)

	return &orderId, &paid, nil
}
