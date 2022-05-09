package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rabigautam/go-serverless-todo/pkg/validators"
)

var (
	ErrorFailedToFetched     = "Error Failed To Fetched"
	ErrorInvalidData         = "Error Invalid Data"
	ErrorFailedToUnmarshal   = "Error Failed To Unmarshal"
	ErrorCouldNotMarshal     = "Error Could Not Marshal"
	ErrorCouldNotDeleteItem  = "Error Could Not Delete Item"
	ErrorCouldNotCreateItem  = "Error Could Not Create Item"
	ErrorUserAlreadyExists   = "Error User Already Exists"
	ErrorUserDoesNotExist    = "Error User Does Not Exist"
	ErrorInvalidEmailAddress = "Error Invalid Email Address"
	ErrorCouldNotUpdateItem  = "Error Could Not Update Item"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	Phone     string `json:"phone"`
	LastName  string `json:"lastname"`
}



func CreateUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	var (
		u User
	)

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmailAddress)
	}
	currentUser, _ := FetchUser(u.Email, tablename, dynaClient)

	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserAlreadyExists)

	}
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorCouldNotMarshal)
	}
	//defining params
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorCouldNotCreateItem)
	}
	return &u, nil
}

func FetchUser(email string, tablename string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}
	result, err := dynaClient.GetItem(input)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorFailedToFetched)
	}
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorFailedToFetched)
	}
	return item, nil

}

func FetchUsers(tablename string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}
	result, err := dynaClient.Scan(input)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorFailedToFetched)
	}
	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorFailedToFetched)
	}

	return items, nil

}
func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	if u.Email == "" {
		return nil, errors.New(ErrorInvalidEmailAddress)
	}
	currentUser, _ := FetchUser(u.Email, tablename, dynaClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(currentUser)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorCouldNotMarshal)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}
	_, err = dynaClient.PutItem(input)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(ErrorCouldNotUpdateItem)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]

	_, err := FetchUser(email, tablename, dynaClient)
	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorUserDoesNotExist)
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}
	_, err = dynaClient.DeleteItem(input)
	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
