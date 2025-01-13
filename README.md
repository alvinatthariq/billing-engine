# Billing Engine

A Go-based billing engine for managing loan schedules, outstanding amounts, and delinquency tracking.

## Features

- Create 50-week loans with specified principal amount and interest rate
- Track weekly payment schedule
- Calculate outstanding balance
- Track delinquent borrowers (missed 2 or more consecutive payments)
- Process weekly payments

## Project Structure

```
billing-engine/
├── cmd/
│   └── billing-engine/     # Main application
├── internal/
│   ├── api/              # REST API handlers
│   └── loan/             # Loan package implementation
└── README.md
```

## API Endpoints

### Create Loan
```
POST /loans
Content-Type: application/json

Request:
{
    "id": "loan123",
    "amount": 5000000,
    "interest_rate": 0.10
}

Response:
{
    "id": "loan123",
    "amount": 5000000,
    "interest_rate": 0.10,
    "weekly_payment": 110000,
    "total_repayment": 5500000
}
```

### Get Outstanding Amount
```
GET /loans/{id}/outstanding

Response:
{
    "outstanding": 5500000
}
```

### Check Delinquency Status
```
GET /loans/{id}/delinquent

Response:
{
    "isDelinquent": false
}
```

### Make Payment
```
POST /loans/{id}/payment
Content-Type: application/json

Request:
{
    "amount": 110000,
    "payment_date": "2025-01-12"
}

Response:
{
    "status": "Payment processed successfully"
}
```

## Usage

To run the server:

```bash
cd cmd/billing-engine
go run main.go
```

The server will start on port 8081.

To run tests:

```bash
go test ./...
```

## Implementation Details

- Loans are created with a principal amount, interest rate, and start date
- Weekly payments are calculated as (principal + interest) / 50 weeks
- Payments must be made for the exact weekly amount
- A borrower becomes delinquent after missing 2 consecutive payments
- Outstanding balance is tracked and updated with each payment
- Payment dates must be provided in "YYYY-MM-DD" format
