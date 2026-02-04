# MercadoPago Payment API

A Go backend API for processing card payments through MercadoPago.

## Setup

1. Copy the environment file:
```bash
cp .env.example
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/main.go
```

## API Endpoints

### Process Payment
**POST** `/api/payments`

Request body:
```json
{
  "amount": 1000,
  "doc_type": "CPF",
  "doc_number": "54960480017",
  "card": {
    "number": "5031433215406351",
    "payment_methods": "master",
    "holder_name": "APRO",
    "expiration_month": "11",
    "expiration_year": "2030",
    "security_code": "123"
  },
  "installments": 1
}
```

Response:
```json
{
  "success": true,
  "message": "Payment processed successfully",
  "data": {
    "payment_id": 123456789,
    "status": "approved",
    "status_detail": "accredited",
    "amount": 1000,
    "installments": 1,
    "payment_method": "master",
    "card_last_four": "6351",
    "created_at": "2024-01-01T12:00:00.000Z"
  }
}
```

### Health Check
**GET** `/health`

Response:
```json
{
  "status": "ok",
  "service": "mercadopago-payment-api"
}
```

## Testing with Postman

1. Import the collection or create a new POST request to `http://localhost:8080/api/payments`
2. Set the Content-Type header to `application/json`
3. Use the example request body above
4. Send the request

## Error Responses

```json
{
  "success": false,
  "error": "Error message description"
}
```

## Dependencies

- Go 1.21+
- Fiber web framework
- MercadoPago Test Account
