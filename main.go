package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CalculationRequest struct {
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type CalculationResponse struct {
	Result string `json:"result"`
}

var displayValue string = "0"
var currentOperand string = ""
var currentOperator string = ""

func mainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleCalculation(w http.ResponseWriter, r *http.Request) {
	var req CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Operation {
	case "C":
		resetCalculator()
	case "+", "-", "*", "/":
		updateDisplayWithOperator(req.Operation)
	case "=":
		computeResult()
	default:
		appendDigit(req.Value)
	}

	resp := CalculationResponse{Result: displayValue}
	json.NewEncoder(w).Encode(resp)
}

func resetCalculator() {
	displayValue = "0"
	currentOperand = ""
	currentOperator = ""
}

func updateDisplayWithOperator(operator string) {
	if currentOperand == "" {
		currentOperand = displayValue
	}
	if currentOperator != "" {
		computeResult()
	}
	displayValue += " " + operator + " "
	currentOperator = operator
}

func computeResult() {
	if currentOperand != "" && currentOperator != "" {
		secondOperand := strings.TrimSpace(displayValue[strings.LastIndex(displayValue, " ")+1:])
		result, err := evaluate(currentOperand, secondOperand, currentOperator)
		if err != nil {
			displayValue = "Error"
		} else {
			displayValue = fmt.Sprintf("%v", result)
			currentOperand = displayValue
		}
		currentOperator = ""
	}
}

func appendDigit(digit string) {
	if displayValue == "0" {
		displayValue = digit
	} else {
		displayValue += digit
	}
}

func evaluate(operand1, operand2, operator string) (float64, error) {
	num1, err1 := strconv.ParseFloat(operand1, 64)
	num2, err2 := strconv.ParseFloat(operand2, 64)
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid operands: %v, %v", operand1, operand2)
	}

	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", operator)
	}
}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/calculate", handleCalculation)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
