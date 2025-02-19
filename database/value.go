package database

import "errors"

type Value struct {
	Data any
	Type string
}

func NewValue(data any, Type string) (Value, error) {
	v := Value{
		Type: Type,
	}

	switch Type {
	case "string":
		if data, ok := data.(string); ok {
			v.Data = data
		} else {
			return Value{}, errors.New("wrong data type")
		}
	case "int":
		if data, ok := data.(int); ok {
			v.Data = data
		} else {
			return Value{}, errors.New("wrong data type")
		}
	case "float64":
		if data, ok := data.(float64); ok {
			v.Data = data
		} else {
			return Value{}, errors.New("wrong data type")
		}
	default:
		return Value{}, errors.New("invalid data type")
	}

	return v, nil
}
