package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
)

func (useCase useCase) CreatePayment(body, requestId string) (*entities.ProcessPaymentRequest, error) {
	order, err := convertToStruct(body, requestId)
	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
		return nil, err
	}

	status := "PENDING"

	order.Status = status

	return useCase.repositoryPayment.CreatePayment(*order, requestId)
}

func (useCase useCase) UpdatePayment(body, requestId string) error {
	payment, err := convertToStruct(body, requestId)
	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
		return err
	}

	orderIdResponse, status, err := useCase.repositoryPayment.UpdatePayment(payment.OrderID, requestId)
	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
		return err
	}

	return useCase.sendSQS(orderIdResponse, status, &requestId)
}

func (useCase useCase) sendSQS(orderID, status, requestId *string) error {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sess)
	queueURL := useCase.queueURL

	orderEvent := entities.ProcessPaymentRequest{
		OrderID: *orderID,
		Status:  *status,
	}

	orderJSON, err := json.Marshal(orderEvent)
	if err != nil {
		fmt.Printf("[RequestId: %p][Error marshaling order request: %v]", requestId, err)
		return err
	}

	messageAttributes := map[string]*sqs.MessageAttributeValue{
		"Source": {
			DataType:    aws.String("String"),
			StringValue: aws.String("payments-events-processor"),
		},
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody:       aws.String(string(orderJSON)),
		QueueUrl:          &queueURL,
		MessageAttributes: messageAttributes,
	})

	if err != nil {
		fmt.Printf("[RequestId: %p][Error sending message to SQS: %v]", requestId, err)
		return err
	}

	return nil
}

func convertToStruct(body, requestId string) (*entities.ProcessPaymentRequest, error) {
	var orderRequest entities.ProcessPaymentRequest
	err := json.Unmarshal([]byte(body), &orderRequest)
	if err != nil {
		fmt.Printf("[RequestId: %s][Error marshaling API Gateway request: %v]", requestId, err)
		return nil, err
	}

	return &orderRequest, nil
}
