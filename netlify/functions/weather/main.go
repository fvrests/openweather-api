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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	latitude := request.QueryStringParameters["latitude"]
	if latitude == "" {
		latitude = "41.8826281"
	}

	longitude := request.QueryStringParameters["longitude"]
	if longitude == "" {
		longitude = "-87.6225291"
	}

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
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            string(bodyBytes),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
