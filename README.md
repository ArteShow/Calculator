This is a simple web calculator.

== Overview ==
This calculator uses the HTTP protocol for communication.
It supports the following operations:

  Addition

  Subtraction

  Multiplication

  Division

  Brackets for order of operations (e.g., 2+22 = 6 and (2+2)(2+2) = 16).

== Requirements ==

  The latest Go version

  A code editor

  Some time

== How to Start ==

  Clone the repository into an empty folder.

  Open the terminal and run:
  go mod tidy

  Now, there should be no more errors in the code.

  Open the terminal again and run:
  go run cmd/main.go

  Wait a little bit.

  A small window with two buttons will appear: Allow and Cancel. Click Allow.

  Congratulations! You just started the calculator.

== Usage ==
  To interact with the calculator, open the Windows terminal:

  Press Win + R, type cmd, and press Enter.

  Important: The backslashes () in the JSON body are required to escape the double quotes ("), ensuring the expression is correctly interpreted by the program.

  Perform a calculation:
  curl -X POST "http://localhost:8082/api/v1/calculate" -H "Content-Type: application/json" -d "{\"expression\": \"2+2\"}"

  Modify {\"expression\": \"2+2\"} to use any arithmetic expression.

  The response will be in the format: {"id": 1}. Note the ID.

  Retrieve all expressions:
  curl -X GET http://localhost:8082/api/v1/expressions

  Example response: {"expression":{"2+2":{"id": 1,"status": 200,"result": 4,"Error":""}}}

  Here, the ID is 1. This will show all previously evaluated expressions.

  Retrieve a specific expression by ID:
  curl -X GET http://localhost:8082/api/v1/expression/{your_id}

  Replace {your_id} with the actual ID from the previous command.

  Example response: {"id": 4,"status": 200,"result": 4,"Error":""}

== Cool Feature ==
  You can add a timer for an arithmetic operation, such as multiplication:

  Open internal/agent.go.

  Go to line 162.

  Modify the timer value in milliseconds:
  TIME_ADDITION_MS = 5 * time.Millisecond

== Need Help? ==
  If you have any issues, feel free to contact me at: sokartemax@gmail.com

