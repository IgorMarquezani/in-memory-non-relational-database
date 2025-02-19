package utils

import "strconv"

func In[T comparable](data T, arr []T) bool {
	for _, v := range arr {
		if v == data {
			return true
		}
	}

	return false
}

func AnyToString(data any) (string, error) {
	var str string

	switch data.(type) {
	case int:
		v := data.(int)
		return strconv.Itoa(v), nil
	case float64:
		v := data.(float64)
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:

	}

	return str, nil
}

func WaitForChan[T any](ch chan T) T {
	var msg T

	for {
		select {
		case msg = <-ch:
			return msg
		}
	}
}
