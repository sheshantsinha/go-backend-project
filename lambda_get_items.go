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
	//"github.com/gorilla/mux"
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
	//min_rating := 0.0
	User := request.QueryStringParameters["email"]

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    if err != nil {
        //fmt.Println("Got error creating session:")
        //fmt.Println(err.Error())
        return events.APIGatewayProxyResponse{
			Body:       "Some error occured",
			StatusCode: 400,
		}, nil
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    // Create the Expression to fill the input struct with.
    // Get all movies in that year; we'll pull out those with a higher rating later
    filt := expression.Name("User").Equal(expression.Value(User))

    // Or we could get by ratings and pull out those with the right year later
    //    filt := expression.Name("info.rating").GreaterThan(expression.Value(min_rating))

    // Get back the title, year, and rating
    //proj := expression.NamesList(expression.Name("title"), expression.Name("year"), expression.Name("info.rating"))

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
        //ProjectionExpression:      expr.Projection(),
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
		//json.NewEncoder(request).Encode(item)
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