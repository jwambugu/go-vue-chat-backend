package util

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

type (
	// MySQL stores mysql config
	MySQL struct {
		DBSource string `yaml:"db_source" mapstructure:"db_source"`
	}

	// DBConfig stores all the database configs
	DBConfig struct {
		MySQL MySQL `yaml:"mysql" mapstructure:"mysql"`
	}

	// Config stores all configuration of the application.
	Config struct {
		AppURL        string   `yaml:"app_url" mapstructure:"app_url"`
		AppEnv        string   `yaml:"app_env" mapstructure:"app_env"`
		AppPort       int      `yaml:"app_port" mapstructure:"app_port"`
		DBConfig      DBConfig `yaml:"db_config" mapstructure:"db_config"`
		EncryptionKey string   `yaml:"encryption_key" mapstructure:"encryption_key"`
		PasetoKey     string   `yaml:"paseto_key" mapstructure:"paseto_key"`
	}
)

// GetAbsolutePath returns the project absolute path from the entry point
func GetAbsolutePath() string {
	_, b, _, _ := runtime.Caller(0)

	basePath := fmt.Sprintf("%s/", filepath.Join(filepath.Dir(b), "../.."))

	return basePath
}

// ReadConfig reads configuration from file or environment variables.
func ReadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config.ReadInConfig:: error loading config - %v", err)
	}

	config := &Config{}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("config.Unmarshal:: error unmarshling config - %v", err)
	}

	return config, nil
}
