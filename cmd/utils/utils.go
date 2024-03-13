package utils

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ConvertStringToInt function to return convert value to string from int
func ConvertStringToInt(value string) int {
	if value == "" {
		value = "0"
	}
	convert, _ := strconv.Atoi(value)
	return convert
}

// ConvertStringToBool function to return convert value to string from bool
func ConvertStringToBool(value string) bool {
	return reflect.ValueOf(value).IsZero()
}

// ConvertStringToTimeSeconds function to return convert value to string from time.Duration
func ConvertStringToTimeSeconds(value string) time.Duration {
	valueInt := ConvertStringToInt(value)
	return time.Duration(valueInt) * time.Second
}

// ConvertStringToTimeMilliseconds function to return convert value to string from time.Duration
func ConvertStringToTimeMilliSeconds(value string) time.Duration {
	valueInt := ConvertStringToInt(value)
	return time.Duration(valueInt) * time.Millisecond
}

// FormatThousandSeparator function to return format with thousand separator for parameter's value
func FormatThousandSeparator(value string) string {
	value = strings.ReplaceAll(value, ".", "") // previous delete '.' or ',' to parameter's value
	re := regexp.MustCompile(`(\d+)(\d{3})`)
	for n := ""; n != value; {
		n = value
		value = re.ReplaceAllString(value, `$1.$2`)
	}
	return value
}

func StringToStringPointer(text string) *string {
	if text == "" {
		return nil
	}

	pointText := new(string)
	*pointText = text
	return pointText
}
