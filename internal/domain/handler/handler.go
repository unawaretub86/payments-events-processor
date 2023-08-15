package handler

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/unawaretub86/payments-events-processor/internal/domain/repository"
	"github.com/unawaretub86/payments-events-processor/internal/domain/repository/database"
	"github.com/unawaretub86/payments-events-processor/internal/domain/usecase"
)

const sqsName = "payments-processor"

func HandleSQSMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	lc, _ := lambdacontext.FromContext(ctx)

	requestId := lc.AwsRequestID

	messageBody, source := getSQSInfo(sqsEvent)

	databaseInstance := createDatabaseInstance()

	repoInstance := repository.NewRepository(databaseInstance)

	useCaseInstance := usecase.NewUsePayment(repoInstance)

	if source == sqsName {
		if err := useCaseInstance.UpdatePayment(messageBody, requestId); err != nil {
			fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
			return err
		}
	}

	_, err := useCaseInstance.CreatePayment(messageBody, requestId)

	return err
}

func getSQSInfo(sqsEvent events.SQSEvent) (string, string) {
	var messageBody string
	var source string

	for _, record := range sqsEvent.Records {
		messageBody = record.Body

		sourceAttr := record.MessageAttributes["Source"]

		source = *sourceAttr.StringValue
	}

	return messageBody, source
}

func createDatabaseInstance() database.Database {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodbClient := dynamodb.New(sess)
	return database.NewDataBase(dynamodbClient)
}
