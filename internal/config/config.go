package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// S3Config holds S3 configuration settings
type S3Config struct {
	Bucket    string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
}

// Init initializes the configuration from config file and environment variables
func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AddConfigPath(".")
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".imgood"))
	}

	viper.SetEnvPrefix("IMGOOD")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	return nil
}

// GetS3Config returns the S3 configuration from viper
func GetS3Config() S3Config {
	return S3Config{
		Bucket:    viper.GetString("s3.bucket"),
		Endpoint:  viper.GetString("s3.endpoint"),
		Region:    viper.GetString("s3.region"),
		AccessKey: viper.GetString("s3.access_key"),
		SecretKey: viper.GetString("s3.secret_key"),
	}
}
