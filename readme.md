# OpenWeather API

[![Netlify Status](https://api.netlify.com/api/v1/badges/0638b0c6-65c2-4112-95e9-d05d0d9335e2/deploy-status)](https://app.netlify.com/sites/fvrests-openweather-api/deploys)

This API provides an endpoint to query the OpenWeather API, primarily used for [Lavender New Tab](https://github.com/fvrests/lavender/). This intermediary API allows the Lavender OpenWeather key to be hidden from the client and applies a rate limit to requests to prevent misuse of the OpenWeather endpoint.

Rate limiting logic based on [rate-limiting without overhead netlify or vercel functions](https://lihbr.com/posts/rate-limiting-without-overhead-netlify-or-vercel-functions).

## Development

_Note: rate limiting non-functional in development environment due to Netlify memory caching implementation_

To develop:
`netlify dev`

To deploy:
`netlify deploy` / `netlify deploy --prod`

## Request parameters

| Parameter | Description                  |
| --------- | ---------------------------- |
| latitude  | Latitude in decimal degrees  |
| longitude | Longitude in decimal degrees |
