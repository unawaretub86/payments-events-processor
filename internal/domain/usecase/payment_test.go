package usecase_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"

	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
	"github.com/unawaretub86/payments-events-processor/internal/domain/repository/mocks"
	"github.com/unawaretub86/payments-events-processor/internal/domain/usecase"
)

func TestCreatePayment(t *testing.T) {
	mockRepo := &mocks.Mocks{}
	mockRepo.CreatePaymentFunc = func(order entities.ProcessPaymentRequest, requestId string) (*entities.ProcessPaymentRequest, error) {
		order.Status = "PENDING"
		return &order, nil
	}

	useCase := usecase.NewUsePayment(mockRepo)

	requestBody := `{"order_id": "1234567890"}`
	requestID := "1234567890"

	payment, err := useCase.CreatePayment(requestBody, requestID)

	if err != nil {
		t.Errorf("Error creating payment: %v", err)
	}

	if payment.OrderID != "1234567890" {
		t.Errorf("Expected orderID to be 1234567890, got %v", payment.OrderID)
	}

	if payment.Status != "PENDING" {
		t.Errorf("Expected status to be PENDING, got %v", payment.Status)
	}
}

func TestCreatePayment_Error(t *testing.T) {
	mockRepo := &mocks.Mocks{}
	useCase := usecase.NewUsePayment(mockRepo)

	invalidBody := `{invalid_field: 1234567890}`
	requestID := "1234567890"

	expectedError := fmt.Errorf("invalid input data")

	mockRepo.CreatePaymentFunc = func(order entities.ProcessPaymentRequest, requestId string) (*entities.ProcessPaymentRequest, error) {
		return nil, expectedError
	}

	result, err := useCase.CreatePayment(invalidBody, requestID)

	assert.Error(t, err, expectedError)
	assert.Nil(t, result)
}

func TestUpdatePayment(t *testing.T) {
	mockRepo := &mocks.Mocks{}
	useCase := usecase.NewUsePayment(mockRepo)

	requestBody := `{"order_id": "1234567890"}`
	requestID := "1234567890"

	mockRepo.UpdatePaymentFunc = func(orderID string, requestId string) (*string, *string, error) {
		orderIdResponse := "updated_order_id"
		status := "UPDATED"
		return &orderIdResponse, &status, nil
	}

	_ = useCase.UpdatePayment(requestBody, requestID)

	mockSQS := mocks.NewMockSQS("us-east-2")
	queueURL := "https://queue.amazonaws.com/80398EXAMPLE/MyQueue"

	_, err := mockSQS.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(`{"order_id": "1234567890"}`),
		QueueUrl:    &queueURL,
	})
	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	if err != nil {
		t.Errorf("Error updating payment: %v", err)
	}
}

func TestUpdatePayment_Error(t *testing.T) {
	mockRepo := &mocks.Mocks{}
	useCase := usecase.NewUsePayment(mockRepo)

	invalidBody := `{invalid_field: 1234567890}`
	requestID := "1234567890"

	expectedError := fmt.Errorf("invalid input data")

	mockRepo.CreatePaymentFunc = func(order entities.ProcessPaymentRequest, requestId string) (*entities.ProcessPaymentRequest, error) {
		return nil, expectedError
	}

	err := useCase.UpdatePayment(invalidBody, requestID)

	assert.Error(t, err, expectedError)
}

func TestUpdatePayment_UpdateError(t *testing.T) {
	// Arrange
	mockRepo := &mocks.Mocks{}
	useCase := usecase.NewUsePayment(mockRepo)

	validBody := `{"order_id": "1234567890"}`
	requestID := "1234567890"

	expectedError := fmt.Errorf("error updating payment")

	mockRepo.UpdatePaymentFunc = func(orderID string, requestId string) (*string, *string, error) {
		return nil, nil, expectedError
	}

	err := useCase.UpdatePayment(validBody, requestID)

	assert.Error(t, err, expectedError)
}
