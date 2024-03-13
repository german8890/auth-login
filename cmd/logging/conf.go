package logging

import (
	"os"
	"strings"

	"autenticacion-ms/cmd/entity"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigVarGlobalLogger() {
	var service entity.Service
	service.Name = os.Getenv("OTEL_SERVICE_NAME")
	if val := os.Getenv("OTEL_RESOURCE_ATTRIBUTES"); val != "" {
		var values = strings.Split(val, ",")
		for _, value := range values {
			if strings.Contains(value, "deployment.environment") {
				service.Environment = strings.Split(value, "=")[1]
				continue
			}
			if strings.Contains(value, "service.version") {
				service.Version = strings.Split(value, "=")[1]
				continue
			}
		}
	}
	//Service = service
}

func ConfigureLogger(loggingLevel string, opts ...zap.Option) (*zap.Logger, error) {
	//zapcore.Encode
	opts = append(opts, zap.AddCallerSkip(1))
	conf := zap.Config{
		Level: configureLevelLogging(loggingLevel),
		//DisableCaller: true,
		DisableStacktrace: true, // TODO: Agregar una variable para habilitar o deshabilitar esta opción
		Development:       false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "duplicate-caller",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		//InitialFields:    pattern,
	}
	//zapcore.NewJSONEncoder(conf)
	return conf.Build(opts...)
}

// NewProductionEncoderConfig returns an opinionated EncoderConfig for
// production environments.
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		//EncodeTime:          zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05:0000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func configureLevelLogging(loggingLevel string) zap.AtomicLevel {
	switch strings.ToLower(loggingLevel) {
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "dpanic":
		return zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}
