package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
)

func (useCase useCase) CreatePayment(body, requestId string) (*string, error) {
	order, err := convertToStruct(body, requestId)
	if err != nil {
		fmt.Printf("[RequestId: %s], [Error: %v]", err, requestId)
		return nil, err
	}

	return useCase.repositoryPayment.CreatePayment(order.OrderID, requestId)
}

func convertToStruct(body, requestId string) (*entities.ProcessPaymentRequest, error) {
	var orderRequest entities.ProcessPaymentRequest
	err := json.Unmarshal([]byte(body), &orderRequest)
	if err != nil {
		fmt.Printf("[RequestId: %s][Error marshaling API Gateway request: %v]", err, requestId)
		return nil, err
	}

	return &orderRequest, nil
}

func (useCase useCase) sendSQS(totalPrice int64, orderID, requestId string) error {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sess)
	queueURL := useCase.queueURL

	orderEvent := entities.ProcessPaymentRequest{
		OrderID: orderID,
		Status:  "HI",
	}

	orderJSON, err := json.Marshal(orderEvent)
	if err != nil {
		fmt.Printf("[RequestId: %s][Error marshaling order request: %v]", err, requestId)
		return err
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(orderJSON)),
		QueueUrl:    &queueURL,
	})

	if err != nil {
		fmt.Printf("[RequestId: %s][Error sending message to SQS: %v]", err, requestId)
		return err
	}

	return nil
}
