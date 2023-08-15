package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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
	status := "PAID"

	update := expression.Set(expression.Name("Status"), expression.Value(status))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
		return nil, nil, err
	}

	primaryKey := map[string]*dynamodb.AttributeValue{
		"OrderId": {
			S: aws.String(orderId),
		},
	}

	if _, err = d.dynamodb.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:                 aws.String(d.table),
		Key:                       primaryKey,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}); err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
		return nil, nil, err
	}

	fmt.Printf("[RequestId: %s], [UpdateItem result: %v]", requestId, orderId)

	return &orderId, &status, nil
}
