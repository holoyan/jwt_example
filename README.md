# Simple JWT token generator server

## Requirements
```bigquery
Go 
MongoDb
```

## Install

```
cp .env.example .env
```
add values in `.env`
```bigquery
go run main.go
```

## Issue Token
```bigquery

curl --location --request POST 'http://localhost:8080/tokens?guid=21563'

```

## Refresh Access/Refresh tokens

```bigquery

curl --location --request POST 'http://localhost:8080/tokens/update' \
--header 'Content-Type: application/json' \
--data-raw '{
  "accessToken": "eyJhbGciOiJIUzI1",
  "refreshToken": "WG5KVUzWg=="
}'

```


