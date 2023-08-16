package mocks

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"

	"github.com/unawaretub86/payments-events-processor/internal/domain/entities"
)

type Mocks struct {
	CreatePaymentFunc func(order entities.ProcessPaymentRequest, requestId string) (*entities.ProcessPaymentRequest, error)
	UpdatePaymentFunc func(orderID string, requestId string) (*string, *string, error)
}

type MockSQS struct {
	sqsiface.SQSAPI
	messages map[string][]*sqs.Message
	region   string
}

func NewMockSQS(region string) *MockSQS {
	return &MockSQS{
		messages: make(map[string][]*sqs.Message),
		region:   "us-east-2",
	}
}

func (m *MockSQS) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if _, ok := m.messages[*in.QueueUrl]; !ok {
		m.messages[*in.QueueUrl] = []*sqs.Message{}
	}

	m.messages[*in.QueueUrl] = append(m.messages[*in.QueueUrl], &sqs.Message{
		Body: in.MessageBody,
	})
	return &sqs.SendMessageOutput{}, nil
}

func (m *Mocks) CreatePayment(order entities.ProcessPaymentRequest, requestId string) (*entities.ProcessPaymentRequest, error) {
	return m.CreatePaymentFunc(order, requestId)
}

func (m *Mocks) UpdatePayment(orderID string, requestId string) (*string, *string, error) {
	return m.UpdatePaymentFunc(orderID, requestId)
}
