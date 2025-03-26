This is a simple web calculator.

## Overview
This calculator uses the HTTP protocol for communication. It supports the following operations:
- Addition
- Subtraction
- Multiplication
- Division
- Brackets for order of operations (e.g., 2+2=4 and (2+2)(2+2)=16).

## Requirements
- The latest Go version
- A code editor
- Some time

## How to Start
1. Clone the repository into an empty folder.
2. Open the terminal and run:
   ```
   go mod tidy
   ```
   Now, there should be no more errors in the code.
3. Open the terminal again and run:
   ```
   go run cmd/main.go
   ```
4. Wait a little bit.
5. A small window with two buttons will appear: Allow and Cancel. Click Allow.
6. Congratulations! You just started the calculator.

## Usage
### Using an outdated version without GUI
To interact with the calculator, open the Windows terminal:
1. Press Win + R, type `cmd`, and press Enter.
2. Important: The backslashes (`\`) in the JSON body are required to escape the double quotes (`"`), ensuring the expression is correctly interpreted by the program.

Perform a calculation:
```
curl -X POST "http://localhost:8082/api/v1/calculate" -H "Content-Type: application/json" -d "{\"expression\": \"2+2\"}"
```
Modify `{\"expression\": \"2+2\"}` to use any arithmetic expression.

The response will be in the format: `{"id": 1}`. Note an ID.

Retrieve all expressions:
```
curl -X GET http://localhost:8082/api/v1/expressions
```
Example response: `{"expression":{"2+2":{"id": 1,"status": 200,"result": 4,"Error":""}}}`

Here, the ID is 1. This will show all previously evaluated expressions.

Retrieve a specific expression by ID:
```
curl -X GET http://localhost:8082/api/v1/expression/{your_id}
```
Replace `{your_id}` with the actual ID from the previous command.

Example response: `{"id": 4,"status": 200,"result": 4,"Error":""}`

### Using the latest version with GUI
To use the graphical interface, follow these steps:
1. Wait till a small window comes out.
2. Use the dropdown menu to select an action:
   - "Calculate" to evaluate an expression.
   - "Expressions" to retrieve all previously evaluated expressions.
   - "Get Expression by ID" to retrieve a specific expression by its ID.
3. Enter the required input in the text field and click the corresponding button to perform an action.
4. The result will be displayed in the text area below.

## Cool Feature
You can add a timer for an arithmetic operation, such as multiplication:
1. Open `internal/agent.go`.
2. Go to line 162.
3. Modify the timer value in milliseconds:
   ```
   TIME_ADDITION_MS = 5 * time.Millisecond
   ```

## Need Help?
If you have any issues, feel free to contact me at: sokartemax@gmail.com