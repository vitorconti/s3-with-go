package configs

import (
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	AwsRegion     string `mapstructure:"AWS_REGION"`
	AwsKey        string `mapstructure:"AWS_KEY"`
	AwsPassword   string `mapstructure:"AWS_PASSWORD"`
	AwsBucketName string `mapstructure:"AWS_BUCKET"`
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {

		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}
