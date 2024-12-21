package calc

import "errors"

var (
	ErrInvalidNumberFormat = errors.New("invalid number format")
	ErrUnsupportedSymbol   = errors.New("unsupported symbol")
	ErrUnbalancedBrackets  = errors.New("unbalanced brackets")
	ErrInvalidExpression   = errors.New("invalid expression")
	ErrNotEnoughOperands   = errors.New("not enough operands")
	ErrDivisionByZero      = errors.New("division by zero")
)
