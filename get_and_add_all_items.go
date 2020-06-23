package main

import (
    "fmt"
    "os"
	"encoding/json"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"net/http"
    "log"
	"github.com/gorilla/mux"
)

// Create structs to hold info about new item
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

// Get the movies with a minimum rating of 8.0 in 2011
func addItem(w http.ResponseWriter, req *http.Request) {
    // Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    av, err := dynamodbattribute.MarshalMap(person)
    if err != nil {
        fmt.Println("Got error marshalling map:")
        fmt.Println(err.Error())
        os.Exit(1)
    }
	//fmt.Println(item.ID)
    // Create item in table Movies
    input := &dynamodb.PutItemInput{
        Item: av,
        TableName: aws.String("UserEvent"),
    }

    _, err = svc.PutItem(input)

    if err != nil {
        fmt.Println("Got error calling PutItem:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println("Successfully added record in UserEvent table")
}

func getItems(w http.ResponseWriter, req *http.Request) {
	//min_rating := 0.0
	paramsx := mux.Vars(req)
    User := paramsx["email"]

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    if err != nil {
        fmt.Println("Got error creating session:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    filt := expression.Name("User").Equal(expression.Value(User))

    expr, err := expression.NewBuilder().WithFilter(filt).Build()

    if err != nil {
        fmt.Println("Got error building expression:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // Build the query input parameters
    params := &dynamodb.ScanInput{
        ExpressionAttributeNames:  expr.Names(),
        ExpressionAttributeValues: expr.Values(),
        FilterExpression:          expr.Filter(),
        TableName:                 aws.String("UserEvent"),
    }

    result, err := svc.Scan(params)

    if err != nil {
        fmt.Println("Query API call failed:")
        fmt.Println((err.Error()))
        os.Exit(1)
    }
	
    for _, i := range result.Items {
        item := Person{}

        err = dynamodbattribute.UnmarshalMap(i, &item)

        if err != nil {
            fmt.Println("Got error unmarshalling:")
            fmt.Println(err.Error())
            os.Exit(1)
        }

        json.NewEncoder(w).Encode(item)
        fmt.Println(item)
        fmt.Println()
    }
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/event/{email}", getItems).Methods("GET")
    router.HandleFunc("/event/", addItem).Methods("POST")
    log.Fatal(http.ListenAndServe(":6000", router))
}