package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Token struct {
	Value string
}

type DeclareOperation struct {
	VariableName string
}

type AssignmentOperation struct {
	VariableName  string
	VariableValue string
}

type IncrementOperation struct {
	VariableName   string
	IncrementValue string
}

type DecrementOperation struct {
	VariableName   string
	DecrementValue string
}

type PrintOperation struct {
	VariableName  string
	VariableValue string
}

type BranchOperation struct {
	CompareTarget      string
	ComparisonOperator string
	Comparison         string
}

type EndOfBranchOperation struct {
}

func main() {
	filePath := os.Args[1]
	scriptData, readError := os.ReadFile(filePath)
	if readError != nil {
		fmt.Println(fmt.Errorf("%v", readError))
		os.Exit(1)
	}

	execute(analyse(parseScript(string(scriptData))))
}

// Reads a script and creates a sequence of string
// tokens based on a set of separators.
func parseScript(script string) []Token {
	var tokenList []Token
	token := ""
	for cursor := 0; cursor < len(script); cursor++ {
		currentRune := script[cursor : cursor+1]
		token = token + currentRune
		if currentRune == " " || currentRune == "\n" {
			cleanToken := strings.Replace(token, "\n", "", -1)
			cleanToken = strings.Replace(cleanToken, " ", "", -1)
			tokenList = append(tokenList, Token{Value: cleanToken})
			token = ""
		}
	}
	return tokenList
}

// Reads through the tokens and creates a sequence
// of logical operations such as printing, declaring
// a variable or setting a value.
func analyse(tokenList []Token) []interface{} {
	var operationList []interface{}
	for cursor := 0; cursor < len(tokenList); cursor++ {
		if tokenList[cursor].Value == "DECLARE" {
			variableName := tokenList[cursor+1].Value
			operationList = append(operationList, DeclareOperation{VariableName: variableName})
			cursor += 1
		} else if tokenList[cursor].Value == "SET" {
			variableName := tokenList[cursor+1].Value
			variableValue := tokenList[cursor+2].Value
			operationList = append(operationList, AssignmentOperation{
				VariableName:  variableName,
				VariableValue: variableValue,
			})
			cursor += 2
		} else if tokenList[cursor].Value == "PRINT" {
			variableName := tokenList[cursor+1].Value
			operationList = append(operationList, PrintOperation{
				VariableName: variableName,
			})
			cursor += 1
		} else if tokenList[cursor].Value == "INCREMENT" {
			variableName := tokenList[cursor+1].Value
			incrementValue := tokenList[cursor+2].Value
			operationList = append(operationList, IncrementOperation{
				VariableName:   variableName,
				IncrementValue: incrementValue,
			})
			cursor += 2
		} else if tokenList[cursor].Value == "DECREMENT" {
			variableName := tokenList[cursor+1].Value
			decrementValue := tokenList[cursor+2].Value
			operationList = append(operationList, DecrementOperation{
				VariableName:   variableName,
				DecrementValue: decrementValue,
			})
			cursor += 2
		} else if tokenList[cursor].Value == "IF" {
			compareTarget := tokenList[cursor+1].Value
			comparisonOperator := tokenList[cursor+2].Value
			comparison := tokenList[cursor+3].Value
			operationList = append(operationList, BranchOperation{
				CompareTarget:      compareTarget,
				ComparisonOperator: comparisonOperator,
				Comparison:         comparison,
			})
			cursor += 3
		} else if tokenList[cursor].Value == "ENDIF" {
			operationList = append(operationList, EndOfBranchOperation{})
		}
	}

	return operationList
}

// Runs through a set of logical operations and
// executes them accordingly, assigning memory
// where necessary.
func execute(operationList []interface{}) {
	var mem = make(map[string]string)
	var shouldSkipProcessing = false
	for _, op := range operationList {
		switch op.(type) {
		case EndOfBranchOperation:
			if shouldSkipProcessing {
				shouldSkipProcessing = false
			}
			break
		case DeclareOperation:
			if shouldSkipProcessing {
				break
			}
			mem[op.(DeclareOperation).VariableName] = "null" // Introduce null reference exception xD
			break
		case AssignmentOperation:
			if shouldSkipProcessing {
				break
			}
			mem[op.(AssignmentOperation).VariableName] = op.(AssignmentOperation).VariableValue
			break
		case PrintOperation:
			if shouldSkipProcessing {
				break
			}
			fmt.Println(mem[op.(PrintOperation).VariableName])
			break
		case IncrementOperation:
			if shouldSkipProcessing {
				break
			}
			initial, _ := strconv.Atoi(mem[op.(IncrementOperation).VariableName])
			valueToAdd, _ := strconv.Atoi(op.(IncrementOperation).IncrementValue)
			result := initial + valueToAdd
			mem[op.(IncrementOperation).VariableName] = fmt.Sprint(result)
		case DecrementOperation:
			if shouldSkipProcessing {
				break
			}
			initial, _ := strconv.Atoi(mem[op.(DecrementOperation).VariableName])
			valueToSubtract, _ := strconv.Atoi(op.(DecrementOperation).DecrementValue)
			result := initial - valueToSubtract
			mem[op.(DecrementOperation).VariableName] = fmt.Sprint(result)
		case BranchOperation:
			if shouldSkipProcessing {
				break
			}

			comparisonTarget, _ := strconv.Atoi(op.(BranchOperation).CompareTarget)
			comparisonOperator := op.(BranchOperation).ComparisonOperator
			comparison, _ := strconv.Atoi(op.(BranchOperation).Comparison)

			match := true
			switch comparisonOperator {
			case "<":
				match = comparisonTarget < comparison
				break
			case ">":
				match = comparisonTarget > comparison
				break
			case "==":
				match = comparisonTarget == comparison
				break
			case "!=":
				match = comparisonTarget != comparison
				break
			}

			if !match {
				shouldSkipProcessing = true
			}
		}
	}
}
