package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/unawaretub86/payments-events-processor/internal/domain/handler"
)

func main() {
	lambda.Start(handler.HandleSQSMessage)
}
