package utils

import (
	"fmt"
	"strconv"
)

func DecimalKeep2(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
