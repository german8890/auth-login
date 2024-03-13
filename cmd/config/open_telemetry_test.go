package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitTracerProvider(t *testing.T) {
	t.Run("Init success Tracer Provider", func(t *testing.T) {
		traceProvider := InitTracerProvider()
		assert.NotNil(t, traceProvider, "Object couldn't be empty")
	})
}
