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

func HandleSQSMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	lc, _ := lambdacontext.FromContext(ctx)

	requestId := lc.AwsRequestID

	var messageBody string
	var source string

	for _, record := range sqsEvent.Records {
		messageBody = record.Body

		sourceAttr := record.MessageAttributes["Source"]

		source = *sourceAttr.StringValue
	}

	databaseInstance := createDatabaseInstance()

	repoInstance := repository.NewRepository(databaseInstance)

	useCaseInstance := usecase.NewUsePayment(repoInstance)

	if source == "payments-processor" {
		if err := useCaseInstance.UpdatePayment(messageBody, requestId); err != nil {
			fmt.Printf("[RequestId: %s], [Error: %v]", requestId, err)
			return err
		}
	}

	_, err := useCaseInstance.CreatePayment(messageBody, requestId)

	return err
}

func createDatabaseInstance() database.Database {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodbClient := dynamodb.New(sess)
	return database.NewDataBase(dynamodbClient)
}
