# API de Busca da Temperatura a partir do CEP

TODO:
- Fazer o docker file
- Fazer o deploy no Google Cloud Run


## Execução
* Localmente
```cmd
    go run main.go
```

* Docker
```cmd
    docker compose up --build
```

* Google Cloud Run
```cmd
    gcloud run deploy --source .
```

## Teste
```cmd
    go test ./...
```

## API
```cmd
    GET /?cep={cep}
```
Retornos
* 200 - `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`
* 422 - `{"error":"invalid zipcode"}` (cep inválido)
* 404 - `{"error":"can not find zipcode"}` (cep não encontrado)


## Cloudrun
https://minha-app-573793284223.southamerica-east1.run.app/?cep=20040001

Deploy
```cmd
gcloud run deploy minha-app --image gcr.io/fullcycle-471609/cloudrun --platform managed --region southamerica-east1 --allow-unauthenticated --set-env-vars WEATHERAPI_KEY=5df680ab2a404916b8895424251209
```
