package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rabigautam/go-serverless-todo/pkg/user"
)

var ErrorMethodNotAllowed = "Method Not Allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}


func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, tableName, dynaClient)
	if err != nil {
		fmt.Println(err)
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)

}
func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynaClient)
		if err != nil {
			fmt.Println(err)
			return apiResponse(http.StatusBadRequest, ErrorBody{
				aws.String(err.Error()),
			})
		}
		return apiResponse(http.StatusOK, result)
	}
	result, err := user.FetchUsers(tableName, dynaClient)
	if err != nil {
		fmt.Println(err)
		return apiResponse(http.StatusInternalServerError, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}


func UnhandledMethod(_ events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)

}
