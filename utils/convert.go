package utils

import (
	"fmt"
	"strconv"
)

func String(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func Float64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	default:
		return 0, fmt.Errorf("Unsupported type: %T", value)
	}
}
