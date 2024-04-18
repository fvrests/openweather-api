package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

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

// Temporarily store request history per IP -- see "rate limiting without overhead"
// Will be lost on cold start, but will be retained in memory for subsequent requests
// info: https://lihbr.com/posts/rate-limiting-without-overhead-netlify-or-vercel-functions
// note: should be tested at deploy link -- memory cache not retained in local environment
var history = make(map[string]time.Time)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// get client IP
	// info: https://answers.netlify.com/t/include-request-ip-address-on-event-object/1820/5
	var ip = request.Headers["x-nf-client-connection-ip"] // user ip

	// previous request from this IP is <15 seconds ago - return / block request
	if history[ip].After(time.Now().Add(-1 * 15 * time.Second)) {
		return &events.APIGatewayProxyResponse{StatusCode: 429, Body: "Rate limit exceeded"}, nil
	}

	// previous request >15 seconds ago - allow request & store the time of this request
	history[ip] = time.Now()

	latitude := request.QueryStringParameters["latitude"]
	longitude := request.QueryStringParameters["longitude"]

	if latitude == "" || longitude == "" {
		return &events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing latitude and/or longitude"}, nil
	}

	// openweather api request
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
