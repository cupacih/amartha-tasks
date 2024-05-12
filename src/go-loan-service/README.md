# Example 3: Loan Service ( system design and abstraction)
we are building a loan engine. A loan can a multiple state: proposed , approved, invested, disbursed. the rule of state:
1. proposed is the initial state (when loan created it will has proposed state):
2. approved is once it approved by our staff. 
   1. approval must contains several information:
      1. the picture proof of the a field validator has visited the borrower 
      2. the employee id of field validator 
      3. date of approval 
   2. once approved it can not go back to proposed state
   3. once approved loan is ready to be offered to investors/lender
3. invested is once total amount of invested is equal the loan principal
   1. loan can have multiple investors, each with each their own amount
   2. total of invested amount can not be bigger than the loan principal amount
   3. once invested all investors will receive an email containing link to agreement letter (pdf)
4. disbursed is when is loan is given to borrower. 
   1. disbursement must contains several information:
      1. the loan agreement letter signed by borrower (pdf/jpeg)
      2. the employee id of the field officer that hands the money and/or collect the agreement letter
      3. date of disbursement


movement between state can only move forward, and a loan only need following information:
- borrower id number
- principal amount
- rate, will define total interest that borrower will pay
- ROI return of investment, will define total profit received by investors
- link to the generated agreement letter


design a RESTFful api that satisfy above requirement.


## Example of Request

GET /loans

```
curl -X GET http://localhost:50052/loans/eccae2a6-9d88-4f08-82be-a80ab235a7e7
```

POST /loans
```
curl -X POST http://localhost:50052/loan \
-H "Content-Type: application/json" \
-d '{
  "borrower_id": "123456789",
  "principal_amount": 1000000,
  "rate": 0.055
}'
```

POST /loan/approve
```

curl -X POST http://localhost:50052/loan/approve \
-H "Content-Type: application/json" \
-d '{
  "loan_id": "eccae2a6-9d88-4f08-82be-a80ab235a7e7",
  "picture_proof_url": "https://example.com/picture.jpg",
  "field_validator_id": "456",
  "approval_date": "2024-05-11T09:58:00Z"
}'
```

POST /loan/invest

```
curl -X POST http://localhost:50052/loan/invest \
-H "Content-Type: application/json" \
-d '{
  "loan_id": "eccae2a6-9d88-4f08-82be-a80ab235a7e7",
  "investors": [
    {
      "investor_id": "123",
      "amount": 600000
    },
    {
      "investor_id": "456",
      "amount": 400000
    }
  ]
}'
```

POST /loan/disburse

```
curl -X POST http://localhost:50052/loan/disburse \
-H "Content-Type: application/json" \
-d '{
  "loan_id": "eccae2a6-9d88-4f08-82be-a80ab235a7e7",
  "signed_agreement_letter": "https://example.com/agreement.pdf",
  "field_officer_id": "123",
  "disbursement_date": "2024-05-11T12:05:00Z"
}'
```