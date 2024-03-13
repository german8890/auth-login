package config_test

import (
	"autenticacion-ms/cmd/config"
	"context"
	"flag"
	"os"
	"testing"

	"autenticacion-ms/cmd/logging"

	"github.com/stretchr/testify/assert"
)

var (
	flagConfig = flag.String("config", "./../../configs/properties_test.yml", "path to the config file")
	logger     = logging.New("*").With(context.Background())
)

func TestLoadFailedYaml(t *testing.T) {
	t.Run("Test Load Fail", func(t *testing.T) {
		_, _ = config.Load(*flagConfig, logger)
		if e := recover(); e != nil {
			assert.NotNil(t, e, "Ok")
		}
	})
}

func TestLoadFailedContextPath(t *testing.T) {
	os.Setenv("LEVEL_LOGGING", "debug")
	os.Setenv("CONTEXT_PATH", "")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("ENABLE_PAYLOAD_LOGGING", "true")
	t.Run("Test Load Fail", func(t *testing.T) {
		_, _ = config.Load(*flagConfig, logger)
		if e := recover(); e != nil {
			assert.NotNil(t, e, "Ok")
		}
	})
}

func TestLoad(t *testing.T) {
	flag.Parse()
	os.Setenv("LEVEL_LOGGING", "debug")
	os.Setenv("CONTEXT_PATH", "/test")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("ENABLE_PAYLOAD_LOGGING", "true")
	t.Run("Test OK", func(t *testing.T) {
		cgf, err := config.Load(*flagConfig, logger)
		assert.Nil(t, err, "An error isn't expected but got no nil.")
		assert.NotNil(t, cgf, "Ok")
	})
}

func TestLoadFailPath(t *testing.T) {
	os.Setenv("LEVEL_LOGGING", "debug")
	os.Setenv("CONTEXT_PATH", "/test")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("ENABLE_PAYLOAD_LOGGING", "true")
	t.Run("Test Load Fail", func(t *testing.T) {
		_, _ = config.Load(*flagConfig, logger)
		if e := recover(); e != nil {
			assert.NotNil(t, e, "Ok")
		}
	})
}

func TestLoadFileNotExistPath(t *testing.T) {
	os.Setenv("LEVEL_LOGGING", "debug")
	os.Setenv("CONTEXT_PATH", "/test")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("ENABLE_PAYLOAD_LOGGING", "true")
	t.Run("Test", func(t *testing.T) {
		_, _ = config.Load("./example.yaml", logger)
		if e := recover(); e != nil {
			assert.NotNil(t, e, "Ok")
		}
	})
}
