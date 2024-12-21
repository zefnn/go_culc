package calc

import (
	"strconv"
	"unicode"
)

// Calc вычисляет значение арифметического выражения
func Calc(expression string) (float64, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}
	return evaluateExpression(tokens)
}

type token struct {
	Type  string
	Value string
}

func tokenize(expression string) ([]token, error) {
	var tokens []token
	for i := 0; i < len(expression); i++ {
		char := expression[i]

		// пропускаем пробелы
		if char == ' ' {
			continue
		}

		// обработка чисел
		if unicode.IsDigit(rune(char)) || char == '.' {
			number := string(char)
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				i++
				number += string(expression[i])
			}
			_, err := strconv.ParseFloat(number, 64)
			if err != nil {
				return nil, ErrInvalidNumberFormat
			}
			tokens = append(tokens, token{Type: "number", Value: number})
			continue
		}

		// обработка операторов и скобок
		switch char {
		case '+', '-', '*', '/', '(', ')':
			tokens = append(tokens, token{Type: "operator", Value: string(char)})
		default:
			return nil, ErrUnsupportedSymbol
		}
	}
	return tokens, nil
}

func evaluateExpression(tokens []token) (float64, error) {
	var numbers []float64
	var operators []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case "number":
			num, _ := strconv.ParseFloat(token.Value, 64)
			numbers = append(numbers, num)

		case "operator":
			if token.Value == "(" {
				operators = append(operators, token.Value)
			} else if token.Value == ")" {
				for len(operators) > 0 && operators[len(operators)-1] != "(" {
					err := applyOperation(&numbers, operators[len(operators)-1])
					if err != nil {
						return 0, err
					}
					operators = operators[:len(operators)-1]
				}
				if len(operators) == 0 {
					return 0, ErrUnbalancedBrackets
				}
				operators = operators[:len(operators)-1] // удаляем открывающую скобку
			} else {
				for len(operators) > 0 && operators[len(operators)-1] != "(" &&
					precedence[operators[len(operators)-1]] >= precedence[token.Value] {
					err := applyOperation(&numbers, operators[len(operators)-1])
					if err != nil {
						return 0, err
					}
					operators = operators[:len(operators)-1]
				}
				operators = append(operators, token.Value)
			}
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return 0, ErrUnbalancedBrackets
		}
		err := applyOperation(&numbers, operators[len(operators)-1])
		if err != nil {
			return 0, err
		}
		operators = operators[:len(operators)-1]
	}

	if len(numbers) != 1 {
		return 0, ErrInvalidExpression
	}

	return numbers[0], nil
}

func applyOperation(numbers *[]float64, operator string) error {
	if len(*numbers) < 2 {
		return ErrNotEnoughOperands
	}

	b := (*numbers)[len(*numbers)-1]
	a := (*numbers)[len(*numbers)-2]
	*numbers = (*numbers)[:len(*numbers)-2]

	var result float64
	switch operator {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return ErrDivisionByZero
		}
		result = a / b
	}

	*numbers = append(*numbers, result)
	return nil
}
