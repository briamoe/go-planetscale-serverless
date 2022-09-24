package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/briamoe/go-planetscale-serverless"
)

type User struct {
	ID   int    `ps:"id" json:"id"`
	Name string `ps:"name" json:"name"`
}

func handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c, err := planetscale.NewConnection(&planetscale.Config{
		Username: os.Getenv("PSCALE_USERNAME"),
		Password: os.Getenv("PSCALE_PASSWORD"),
		Host:     os.Getenv("PSCALE_HOST"),
	})
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "unable to connect to planetscale",
		}, nil
	}

	e, err := c.Execute("SELECT * FROM users WHERE id=?", req.QueryStringParameters["id"])
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "unable to execute statement",
		}, nil
	}

	var u []*User
	if err := e.Decode(&u); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "unable to decode users",
		}, nil
	}

	m, err := json.Marshal(u)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "unable to marshal users",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(m),
	}, nil
}

func main() {
	lambda.Start(handle)
}
