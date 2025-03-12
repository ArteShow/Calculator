package calculate

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Create Logger
func createLogger(folderPath, fileName string) (*log.Logger, error, int) {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return nil, errors.New("Internal server error"), 500
	}

	filePath := folderPath + "/" + fileName

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.New("Internal server error"), 500
	}

	logger := log.New(file, "", log.LstdFlags)
	return logger, nil, 200
}

func Postpone(nums []float64, index int) []float64 {
	return append(nums[:index], nums[index+1:]...)
}

func PostponeStringSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

// Calculate the expression without brackets
func CalcBasic(expression string) (float64, error, int) {
	// Setup the logger
	Logger, err, status := createLogger("../log", "CalcBasicLog.txt")
	if err != nil {
		Logger.Println("[ERROR]:", err)
		return 0.0, err, status
	}
	Logger.Println("///////////////Calculation started/////////////")

	// Some arrays and variables
	var RuneWithout string
	MathOperators := make([]string, 0)
	Numbers := make([]float64, 0)
	var num string
	expression = strings.ReplaceAll(expression, ",", ".")
	for _, letter := range expression {
		if letter == ' ' {
			continue
		}
		RuneWithout += string(letter)
	}
	Logger.Println("Removed all spaces and replaced all commas with periods")

	String := string(RuneWithout)
	if len(String) == 0 {
		Logger.Println("[ERROR]: idk")
		return 0.0, errors.New("Something went wrong. Please try again"), 404
	}

	// Checking for letters
	for _, letter := range String {
		if unicode.IsDigit(letter) || letter == '.' {
			num += string(letter)
		} else if unicode.IsLetter(letter) {
			return 0.0, errors.New("There is a letter in the expression"), 422
		} else {
			if num != "" {
				numb, err := strconv.ParseFloat(num, 64)
				if err != nil {
					return 0.0, err, 404
				}
				Numbers = append(Numbers, numb)
				num = ""
			}
		}
	}

	if num != "" {
		numb, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0.0, err, 404
		}
		Numbers = append(Numbers, numb)
	}

	// Count the mathematical symbols
	Logger.Println("Cleaned up the expression")
	var lastOperator string
	for _, letter := range String {
		newLetter := string(letter)
		if newLetter == "+" || newLetter == "-" || newLetter == "*" || newLetter == "/" {
			if lastOperator == newLetter {
				Logger.Println("[ERROR]: error while checking the mathematical symbols")
				return 0.0, errors.New("Something went wrong. Please try without minus"), 422
			}
			lastOperator = newLetter
			MathOperators = append(MathOperators, newLetter)
		}
	}

	Logger.Println("Counted all the mathematical symbols.")
	for i := 0; i < len(MathOperators); i++ {
		if MathOperators[i] == "*" || MathOperators[i] == "/" {
			if MathOperators[i] == "*" {
				if i+1 < len(Numbers) {
					Numbers[i] = Numbers[i] * Numbers[i+1]
				} else {
					Logger.Println("[ERROR]: error by multiplication")
					return 0.0, errors.New("Multiply error"), 404
				}
			} else if MathOperators[i] == "/" {
				if Numbers[i+1] == 0 {
					Logger.Println("[ERROR]: error divide by zero")
					return 0.0, errors.New("Division by zero"), 422
				}
				if i+1 < len(Numbers) {
					Numbers[i] = Numbers[i] / Numbers[i+1]
				} else {
					Logger.Println("[ERROR]: error by division")
					return 0.0, errors.New("Error by division"), 404
				}
			}
			Numbers = Postpone(Numbers, i+1)
			MathOperators = PostponeStringSlice(MathOperators, i)
			i--
		}
	}

	// Do the addition and subtraction
	Logger.Println("Performed division and multiplication")
	for i := 0; i < len(MathOperators); i++ {
		if MathOperators[i] == "-" {
			Numbers[i] = Numbers[i] - Numbers[i+1]
			Numbers = Postpone(Numbers, i+1)
			MathOperators = PostponeStringSlice(MathOperators, i)
			i--
		} else if MathOperators[i] == "+" {
			Numbers[i] = Numbers[i] + Numbers[i+1]
			Numbers = Postpone(Numbers, i+1)
			MathOperators = PostponeStringSlice(MathOperators, i)
			i--
		}
	}
	Logger.Println("Performed addition and subtraction")
	if len(Numbers) == 0 {
		Logger.Println("[ERROR]: error by plus and minus")
		return 0.0, errors.New("Error by calculation"), 404
	}
	Logger.Println("End of the function CalcBasic")
	return Numbers[0], nil, 200
}

// Calculation with brackets
func Calc(expression string) (float64, error, int) {
	// Setup the logger
	Logger, err, status := createLogger("../log", "CalcLog.txt")
	if err != nil {
		Logger.Println("[ERROR]:", err)
		return 0.0, err, status
	}
	Logger.Println("/////////////Check for Brackets///////////////")

	// Preparing the expression
	expression = strings.ReplaceAll(expression, ",", ".")
	var ResultString string
	stack := []string{}

	// Delete all spaces
	for _, letter := range expression {
		if letter == ' ' {
			continue
		}
		ResultString += string(letter)
	}
	Logger.Println("Replaced all commas with periods")

	// Checking the brackets
	ResultString = "(" + ResultString + ")"
	for i := 0; i < len(ResultString); i++ {
		if ResultString[i] == '(' {
			stack = append(stack, ResultString[i:i+1])
		} else if ResultString[i] == ')' {
			if len(stack) == 0 {
				Logger.Println("[ERROR]: error by counting the brackets")
				return 0.0, errors.New("Error by counting the brackets"), 422
			}
			Logger.Println("Checked for wrong brackets")
			stack = stack[:len(stack)-1]
			startIndex := strings.LastIndex(ResultString[:i], "(")
			innerExpression := ResultString[startIndex+1 : i]
			Result, err, code := CalcBasic(innerExpression)
			if err != nil {
				return 0.0, err, code
			}
			ResultString = ResultString[:startIndex] + strconv.FormatFloat(Result, 'f', 2, 64) + ResultString[i+1:]
			i = startIndex + len(strconv.FormatFloat(Result, 'f', 2, 64)) - 1
		}
	}

	// Check the calculation
	Logger.Println("Calculated all expressions step by step")
	if len(stack) != 0 {
		Logger.Println("[ERROR]: too many ( )")
		return 0.0, errors.New("Too many brackets"), 422
	}
	result, err, code := CalcBasic(ResultString)
	if err != nil {
		Logger.Println("[ERROR]: error in CalcBasic")
		return 0.0, err, code
	}
	Logger.Println("End of the function Calc")
	return result, nil, 200
}
