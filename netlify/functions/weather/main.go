package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var HEADERS = map[string]string{
	"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin",
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Methods": "POST, OPTIONS",
	"Access-Control-Max-Age":       "8640",
	"Access-Control-Allow-Origin":  "*",
	"Vary":                         "Origin",
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	latitude := request.QueryStringParameters["latitude"]

	longitude := request.QueryStringParameters["longitude"]

	if latitude == "" || longitude == "" {
		log.Println("latitude and longitude are required")
	}

	baseUrl := "https://api.openweathermap.org/data/2.5/weather"
	// local .env for dev environment, netlify env for production
	appId := os.Getenv("OPENWEATHER_KEY")
	url := baseUrl + "?lat=" + latitude + "&lon=" + longitude + "&appid=" + appId

	res, err := http.Get(url)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         HEADERS,
		Body:            string(bodyBytes),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
