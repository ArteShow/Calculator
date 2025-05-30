# Simple Calculator

## Overview
This calculator uses the HTTP protocol for communication. It supports login and the following operations:
- **Addition**
- **Subtraction**
- **Multiplication**
- **Division**
- **Brackets** for order of operations (e.g., `2+2=4` and `(2+2)(2+2)=16`).

---

## Requirements
- The latest Go version
- A code editor
- Some time

---

## How to Start
1. Clone the repository into an empty folder.
2. Open the terminal and run:
    go mod tidy
    Now, there should be no more errors in the code. 
3. Open the terminal again and run:
    go run cmd/main.go
    Wait a little bit.
4. A small window with two buttons will     appear: Allow and Cancel. Click Allow.

Congratulations! You just started the calculator.

## Usage
To interact with the calculator, open the Windows terminal:

Press Win + R, type cmd, and press Enter.
Important: The backslashes (\) in the JSON body are required to escape the double quotes ("), ensuring the expression is correctly interpreted by the program.

## Register into the program:
    curl -X POST http://localhost:8082/api/v1/register -H "Content-Type: application/json" -d "{\"login\": \"Your username\", \"passwort\": \"Your passwort\"}"

## Login into the program:
    curl -X POST http://localhost:8082/api/v1/login -H "Content-Type: application/json" -d "{\"login\": \"Your username\", \"passwort\": \"Your passwort\"}"

As a response, you will get the JWT key. Note that you will need it later.

## Perform a calculation:
    curl -X POST http://localhost:8082/api/v1/calculate -H "Content-Type: application/json" -H "Authorization: Bearer (your token)" -d "{\"expression\": \"(your expression)\"}"
The response will be in the format:
    {"id": 1}

Note the ID.

## Retrieve all expressions:
    curl -X GET http://localhost:8082/api/v1/expressions -H "Authorization: Bearer (your token)"

Example response:
    {
        "expression": {
            "2+2": {
                "id": 1,
                "status": 200,
                "result": 4,
                "Error": ""
            }
        }
    }

This will show all previously evaluated expressions.

## Retrieve a specific expression by ID:

curl -X GET http://localhost:8082/api/v1/expression/{your_id} -H "Authorization: Bearer (your token)"

Replace {your_id} with the actual ID from the previous command.

Example response:
    {
        "id": 4,
        "status": 200,
        "result": 4,
        "Error": ""
    }

## Need Help?
If you have any issues, feel free to contact me at: sokartemax@gmail.com