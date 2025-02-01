package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"net/url"
	"time"
)

type Config struct {
	App        `yaml:"app"`
	HTTPServer `yaml:"http-server"`
	DB         `yaml:"database"`
}

type App struct {
	Name    string `yaml:"name"`
	Profile string
}

type HTTPServer struct {
	Port string `mapstructure:"port"`
}

type DB struct {
	Host         string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port         uint16        `yaml:"port" env:"DB_PORT" env-default:"port"`
	User         string        `yaml:"user" env:"DB_USER" env-default:"happy_miner"`
	Password     string        `yaml:"password" env:"DB_PASSWORD" env-default:"happy_miner"`
	Name         string        `yaml:"name" env:"DB_NAME" env-default:"server_manager"`
	MaxPoolSiz   int           `yaml:"max-pool-size" env-default:"2"`
	ConnAttempts int           `yaml:"conn-attempts" env-default:"3"`
	ConnTimeout  time.Duration `yaml:"conn-timeout" env-default:"30s"`
}

func New() *Config {
	var cfg Config
	var env string
	flag.StringVar(&env, "env", "local", "Application profile")
	flag.Parse()

	if err := cleanenv.ReadConfig(fmt.Sprintf("./config/config.%s.yaml", env), &cfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &cfg
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
