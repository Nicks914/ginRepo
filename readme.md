## Installation and Setup
    Clone the Repository:

    bash
    Copy code
    git clone https://github.com/Nicks914/ginRepo.git
    cd ginRepo

## Database Setup:

    Create a MySQL database named cetec.
    Import the schema and initial data using the provided SQL scripts (schema.sql and data.sql).
    Environment Variables:

    Open main.go and update the database connection string in the init() function (db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/cetec")).
    Run the Application:

    bash
    Copy code
    go run main.go
    The application will start on http://localhost:8080.

## Endpoints
    1. GET /person/{person_id}/info
    Retrieve person information by ID.

    Example: GET /person/1/info

    Response:

    json
    Copy code
    {
    "name": "mike",
    "phone_number": "444-444-4444",
    "city": "Austin",
    "state": "TX",
    "street1": "213 South 1st St",
    "street2": "",
    "zip_code": "78704"
    }
    2. POST /person/create
    Create a new person record.

    Request Body:

    json
    Copy code
    {
    "name": "YOURNAME",
    "phone_number": "123-456-7890",
    "city": "Sacramento",
    "state": "CA",
    "street1": "112 Main St",
    "street2": "Apt 12",
    "zip_code": "12345"
    }
    Example:

    bash
    Copy code
    curl -X POST \
    http://localhost:8080/person/create \
    -H 'Content-Type: application/json' \
    -d '{
            "name": "John Doe",
            "phone_number": "555-123-4567",
            "city": "Los Angeles",
            "state": "CA",
            "street1": "456 Elm St",
            "street2": "Unit 3",
            "zip_code": "90001"
        }'
    Response:

    json
    Copy code
    {
    "message": "Person created successfully"
    }
## Error Handling
    The application handles errors gracefully and responds with appropriate HTTP status codes (400, 500) and error messages in JSON format.