package parser

import (
	"app/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	simple    = []byte{'+', '-', ':', '_', '#', ',', '('}
	aggregate = []byte{'$', '*', '!', '=', '%', '`', '~', '>'}
)

func ParseVariable(variable []byte) (any, string, error) {
	if len(variable) < 1 {
		return nil, "", errors.New("Invalid token")
	}

	if utils.In(variable[0], simple) {
		return ParseSimple(variable)
	}
	if utils.In(variable[0], aggregate) {
		return ParseAggregate(variable)
	}

	return nil, "", errors.New("Invalid token")
}

func ParseSimple(variable []byte) (any, string, error) {
	var (
		v any
		t string
	)

	tokens := strings.Split(string(variable), "\\r\\n")

	if len(tokens) < 1 {
		return nil, "", errors.New("invalid data type")
	}

	if len(tokens[0]) < 2 {
		return nil, "", errors.New("syntax error")
	}

	switch tokens[0][0] {
	case '+':
		v = tokens[0][1:]
		t = "string"
	case '-':

	case ':':
		data, err := strconv.Atoi(tokens[0][1:])
		if err != nil {
			return nil, "", errors.New("syntax error")
		}

		v = data
		t = "int"
	case ',':
		data, err := strconv.ParseFloat(tokens[0][1:], 64)
		if err == nil {
			return nil, "", errors.New("syntax error")
		}

		v = data
		t = "float64"
	default:
		return nil, "", errors.New("Err: Invalid data type: " + string(tokens[0][0]))
	}

	return v, t, nil
}

func ParseAggregate(variable []byte) (any, string, error) {
	var (
		v any
		t string
	)

	tokens := strings.Split(string(variable), "\\r\\n")

	fmt.Println(tokens)

	if len(tokens) < 2 {
		return nil, "", errors.New("invalid data type")
	}

	v = tokens[1]

	switch tokens[0][0] {
	case '$':
		t = "string"
	default:
		return nil, "", errors.New("Err: Invalid data type: " + string(tokens[0][0]))
	}

	return v, t, nil
}
