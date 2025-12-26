package sql

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func Parse(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val), nil
	case time.Time:
		return fmt.Sprintf("\"%s\"", v.Format(time.DateOnly)), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int:
		return strconv.Itoa(v), nil
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', 2, 64), nil

	}
	return "", errors.New("cannot parse given value to string")
}
