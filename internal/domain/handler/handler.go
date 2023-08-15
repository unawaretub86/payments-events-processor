package handler

import (
	"context"

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

	for _, record := range sqsEvent.Records {
		messageBody = record.Body
	}

	databaseInstance := createDatabaseInstance()

	repoInstance := repository.NewRepository(databaseInstance)

	useCaseInstance := usecase.NewUseOrder(repoInstance)

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
