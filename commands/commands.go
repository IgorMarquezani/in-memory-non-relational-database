package commands

import (
	"app/instruction"
	"app/parser"
	"bytes"
	"errors"
	"strings"
)

func Echo(args []string) (string, error) {
	if len(args) < 1 || len(args) > 1 {
		return "", errors.New("(error) ERR wrong number of arguments for 'echo' command")
	}

	msg, _, err := parser.ParseVariable([]byte(args[0]))
	if err != nil {
		return "", errors.New("(error) ERR " + err.Error())
	}

	v, ok := msg.(string)
	if !ok {
		return "", errors.New("(error) ERR unexpected error")
	}

	return v, nil
}

func Set(args []string, queue chan instruction.Instruction) (chan instruction.Message, error) {
	if len(args) < 2 || len(args) > 5 {
		return nil, errors.New("(error) ERR wrong number of arguments for 'set' command")
	}

	arr := bytes.Trim([]byte(args[1]), "\x00")

	data, Type, err := parser.ParseVariable(arr)
	if err != nil {
		return nil, errors.New("(error) ERR " + err.Error())
	}

	cmd := instruction.Instruction{
		Command: "set",
		Key:     args[0],
		Data:    data,
		Type:    Type,
		Channel: make(chan instruction.Message, 1),
	}

	queue <- cmd

	return cmd.Channel, nil
}

func Get(args []string, queue chan instruction.Instruction) (chan instruction.Message, error) {
	str := strings.Trim(args[0], "\x00")

	cmd := instruction.Instruction{
		Command: "get",
		Key:     str,
		Channel: make(chan instruction.Message, 1),
	}

	queue <- cmd

	return cmd.Channel, nil
}
