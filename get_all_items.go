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
	fmt.Printf("marshalled struct: %+v", av)
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

    fmt.Println("Successfully added 'The Big New Movie' (2015) to Movies table")
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

    // Create the Expression to fill the input struct with.
    // Get all movies in that year; we'll pull out those with a higher rating later
    filt := expression.Name("User").Equal(expression.Value(User))

    // Or we could get by ratings and pull out those with the right year later
    //    filt := expression.Name("info.rating").GreaterThan(expression.Value(min_rating))

    // Get back the title, year, and rating
    //proj := expression.NamesList(expression.Name("title"), expression.Name("year"), expression.Name("info.rating"))

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
        //ProjectionExpression:      expr.Projection(),
        TableName:                 aws.String("UserEvent"),
    }

    // Make the DynamoDB Query API call
    result, err := svc.Scan(params)

    if err != nil {
        fmt.Println("Query API call failed:")
        fmt.Println((err.Error()))
        os.Exit(1)
    }

    num_items := 0
	
    for _, i := range result.Items {
        item := Person{}

        err = dynamodbattribute.UnmarshalMap(i, &item)

        if err != nil {
            fmt.Println("Got error unmarshalling:")
            fmt.Println(err.Error())
            os.Exit(1)
        }

        // Which ones had a higher rating?
        // if item.Info.Rating > min_rating {
        //     // Or it we had filtered by rating previously:
        //     //   if item.Year == year {
        //     num_items += 1
			json.NewEncoder(w).Encode(item)
            fmt.Println(item)
            fmt.Println()
        //}
    }
	//json.NewEncoder(w).Encode(&Person{})
    fmt.Println("Found", num_items, "movie(s) with a rating above", User)
}

func main() {
    router := mux.NewRouter()
    //people = append(people, Person{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}})
    //people = append(people, Person{ID: "2", Firstname: "Maria", Lastname: "Raboy"})
    //router.HandleFunc("/event", GetPeopleEndpoint).Methods("GET")
    router.HandleFunc("/event/{email}", getItems).Methods("GET")
    router.HandleFunc("/event/", addItem).Methods("POST")
    // router.HandleFunc("/event/{id}", UpdatePersonEndpoint).Methods("PATCH")
    log.Fatal(http.ListenAndServe(":12345", router))
}