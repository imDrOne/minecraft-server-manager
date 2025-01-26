package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"net/url"
)

const (
	ProfileLocal = "local"
	ProfileDev   = "dev"
)

type Config struct {
	App        `mapstructure:"app"`
	HTTPServer `mapstructure:"http-server"`
	DB         `mapstructure:"database"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Profile string
}

type HTTPServer struct {
	Port string `mapstructure:"port"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func (d *DB) BuildConnectionString(sslMode string, additionalParams map[string]string) string {
	user := url.QueryEscape(d.User)
	password := url.QueryEscape(d.Password)
	host := d.Host
	port := d.Port
	database := url.QueryEscape(d.Name)

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)

	params := url.Values{}
	params.Add("sslmode", sslMode)

	for key, value := range additionalParams {
		params.Add(key, value)
	}

	if len(params) > 0 {
		connectionString += "?" + params.Encode()
	}

	return connectionString
}

func New() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("app.profile", ProfileLocal)

	if err := viper.BindEnv("app.profile", "APP_PROFILE"); err != nil {
		panic("can't bind app.profile")
	}

	profile := viper.GetString("app.profile")
	slog.Info(fmt.Sprintf("Active profile: %s", profile))

	viper.SetConfigName(fmt.Sprintf("config.%s", profile))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	return config
}
