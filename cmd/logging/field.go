package logging

import "go.uber.org/zap"

type Field struct {
	Key   string
	value interface{}
}

func AnyField(Key string, value interface{}) Field {
	return Field{Key, value}
}

func toFieldZap(fields ...Field) []zap.Field {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.value))
	}
	return zapFields
}
