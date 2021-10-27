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

	// DBConfig stores all the db configs
	DBConfig struct {
		MySQL MySQL `yaml:"mysql"`
	}

	// Config stores all configuration of the application.
	Config struct {
		AppURL        string   `yaml:"app_url" mapstructure:"db_source"`
		AppEnv        string   `yaml:"app_env" mapstructure:"db_source"`
		AppPort       int      `yaml:"app_port" mapstructure:"db_source"`
		DBConfig      DBConfig `yaml:"db_config" mapstructure:"db_source"`
		EncryptionKey string   `yaml:"encryption_key" mapstructure:"db_source"`
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
