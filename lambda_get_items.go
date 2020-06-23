package main

import (
	"encoding/json"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Person struct {
    ID        int   `json:"ID,omitempty"`
    Name string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Status  string   `json:"status,omitempty"`
	Schedule   Schedule `json:"schedule,omitempty"`
	User string `json:"User,omitempty"`
}

type Schedule struct {
    Start  string `json:"start_time,omitempty"`
    Stop string `json:"stop_time,omitempty"`
}

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
func getItems(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	User := request.QueryStringParameters["email"]

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    if err != nil {
        return events.APIGatewayProxyResponse{
			Body:       "Some error occured",
			StatusCode: 400,
		}, nil
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    filt := expression.Name("User").Equal(expression.Value(User))

    expr, err := expression.NewBuilder().WithFilter(filt).Build()

    if err != nil {
        return events.APIGatewayProxyResponse{
			Body:       "Some error occured",
			StatusCode: 400,
		}, nil
    }

    // Build the query input parameters
    params := &dynamodb.ScanInput{
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        FilterExpression:          expr.Filter(),
        TableName:                 aws.String("UserEvent"),
    }

    // Make the DynamoDB Query API call
    result, err := svc.Scan(params)

    if err != nil {
        return events.APIGatewayProxyResponse{
			Body:  "Some error occured",
			StatusCode: 400,
		}, nil
    }
	
    for _, i := range result.Items {
        item := Person{}

        err = dynamodbattribute.UnmarshalMap(i, &item)

        if err != nil {
            return events.APIGatewayProxyResponse{
				Body: "Some error occured",
				StatusCode: 400,
			}, nil
		}
	}
	b, err := json.Marshal(result.Items)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Some error occured",
			StatusCode: 400,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:  string(b),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin" : "*",
		},
	}, nil
}

func main() {
    lambda.Start(getItems)
}