package logging

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type duplicateCallerEncoder struct {
	zapcore.Encoder
}

func (e *duplicateCallerEncoder) Clone() zapcore.Encoder {
	return &duplicateCallerEncoder{Encoder: e.Encoder.Clone()}
}

func (e *duplicateCallerEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	fields = deleteRepeatByLastKey(fields)
	return e.Encoder.EncodeEntry(entry, fields)
}

func deleteRepeatByLastKey(arr []zapcore.Field) (res []zapcore.Field) {
	visited := map[string]bool{}
	for i := len(arr) - 1; i >= 0; i-- {
		n := arr[i].Key
		if visited[n] {
			continue
		}

		visited[n] = true
		res = append([]zapcore.Field{arr[i]}, res...)
	}

	return
}
