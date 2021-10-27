//+build integration

package util

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		expectsError bool
	}{
		{
			name:         "can read config from file",
			path:         GetAbsolutePath(),
			expectsError: false,
		},
		{
			name:         "fails to read config due to invalid file path",
			path:         "",
			expectsError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer viper.Reset()

			config, err := ReadConfig(testCase.path)

			if testCase.expectsError {
				assert.Error(t, err)
				assert.Nil(t, config)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, config)
			assert.IsType(t, &Config{}, config)
		})
	}
}
