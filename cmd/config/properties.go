package config

import (
	"autenticacion-ms/cmd/logging"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/qiangxue/go-env"
	"gopkg.in/yaml.v2"
)

const (
	defaultHttpServerPort = "8080"
	defaultLevelLogging   = "*"
	enablePayloadLogging  = false
)

type Config struct {
	// the server port. Defaults to 8080
	HttpServerPort string `yaml:"server_http_port"`
	// the server port. Defaults to 50001
	GrpcServerPort string `yaml:"server_grpc_port"`
	// Context Path of API Rest
	ContextPath string `yaml:"context_path"`
	// Logger Level Mode e.g: info, error, warn, debug, *
	LevelLogging         string `yaml:"level_logging" env:"LEVEL_LOGGING"`
	EnablePayloadLogging bool
	DocumentsTypeAllowed string `yaml:"documents_type_allowed"`
	Redis                Redis  `yaml:"redis"`
}

type Redis struct {
	RedisAddr     string `yaml:"redis_addr"`
	RedisPassword string `yaml:":redis_pass"`
}

// Validate 'validates' the application configuration
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ContextPath, validation.Required))
}

// Load returns an application configuration which is populated from the given configuration file and environment variables

func Load(file string, logger logging.Logger) (*Config, error) {

	// default config
	c := Config{
		HttpServerPort:       defaultHttpServerPort,
		LevelLogging:         defaultLevelLogging,
		EnablePayloadLogging: false,
	}

	// load from YAML config file
	//bytes, err := ioutil.ReadFile(file)
	file = filepath.Clean(file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	//Load environment variables - LOCAL
	_ = godotenv.Load("./.env")

	// replacing values from ENV
	replaceValuesEnvInFile(&bytes)
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// load from environment variables prefixed with "APP_"
	_ = env.New("APP_", nil).Load(&c)

	// validation
	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}

// replaceValuesEnvInFile, replace value env in file that regex ${&ENV_value}
func replaceValuesEnvInFile(output *[]byte) *[]byte {
	for _, line := range strings.Split(string(*output), "\n") {
		if ok := strings.Contains(line, "${"); ok {
			test := strings.Split(line, "${")
			for i := 1; i < len(test); i++ {
				if strings.Contains(test[i], ":") {
					test[i] = strings.ReplaceAll(test[i], ":", "")
				}
				replaceValueToEnv(test[i], output)
			}
		}

	}
	return output
}

func replaceValueToEnv(line string, output *[]byte) {
	valueToReplace, _, _ := strings.Cut(line, "}")
	if value, exist := os.LookupEnv(valueToReplace); exist {
		valueToReplace := fmt.Sprintf("${%v}", valueToReplace)
		*output = bytes.Replace(*output, []byte(valueToReplace), []byte(value), -1)
	}
}
