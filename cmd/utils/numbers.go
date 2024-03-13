package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const decimalPrecision = "%.10f"

func FloatToString(value float64, precision string) string {
	return fmt.Sprintf(precision, value)
}

func IntToString(value int, precision string) string {
	valueAsString := strconv.Itoa(value)
	return valueAsString
}

func FloatToJSON(value float64) *json.Number {
	valueAsString := FloatToString(value, decimalPrecision)
	valueAsJSON := json.Number(valueAsString)
	return &valueAsJSON
}

func IntToJSON(value int) *json.Number {
	valueAsString := IntToString(value, decimalPrecision)
	valueAsJSON := json.Number(valueAsString)
	return &valueAsJSON
}

func JSONToFloat(value *json.Number) float64 {
	floatValue, _ := value.Float64()
	return floatValue
}
