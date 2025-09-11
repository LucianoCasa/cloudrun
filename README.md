# CEP Weather API

Simple HTTP service written in Go that receives a Brazilian CEP, resolves the
city using ViaCEP and then fetches the current temperature from WeatherAPI.
It returns the temperature in Celsius, Fahrenheit and Kelvin.

## Running locally

```bash
go test ./...

WEATHERAPI_KEY=your_key go run main.go
```

Or using docker-compose:

```bash
WEATHERAPI_KEY=your_key docker compose up --build
```

## API

```
GET /?cep=<8 digit CEP>
```

### Success (200)
```json
{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

### Invalid CEP (422)
```
invalid zipcode
```

### CEP not found (404)
```
can not find zipcode
```

## Deploy to Cloud Run

Build and push the image then deploy:

```bash
gcloud builds submit --tag gcr.io/PROJECT_ID/cep-weather

gcloud run deploy cep-weather \
  --image gcr.io/PROJECT_ID/cep-weather \
  --platform managed \
  --allow-unauthenticated \
  --region us-central1 \
  --set-env-vars WEATHERAPI_KEY=$WEATHERAPI_KEY
```

Replace `PROJECT_ID` with your Google Cloud project and set the
`WEATHERAPI_KEY` environment variable.
