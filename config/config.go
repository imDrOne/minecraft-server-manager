package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	App        `yaml:"app"`
	HTTPServer `yaml:"http-server"`
	DB         `yaml:"database"`
	SSHKeygen  `yaml:"ssh-keygen"`
}

type App struct {
	Name    string `yaml:"name"`
	Profile string
}

type HTTPServer struct {
	Port string `yaml:"port"`
}

type DB struct {
	Host         string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port         string        `yaml:"port" env:"DB_PORT" env-default:"port"`
	User         string        `yaml:"user" env:"DB_USER" env-default:"happy_miner"`
	Password     string        `yaml:"password" env:"DB_PASSWORD" env-default:"happy_miner"`
	Name         string        `yaml:"name" env:"DB_NAME" env-default:"server_manager"`
	MaxPoolSiz   int           `yaml:"max-pool-size" env-default:"2"`
	ConnAttempts int           `yaml:"conn-attempts" env-default:"3"`
	ConnTimeout  time.Duration `yaml:"conn-timeout" env-default:"30s"`
}

type SSHKeygen struct {
	Bits       int    `yaml:"bits" env-default:"2048"`
	Passphrase string `yaml:"passphrase" env:"SSH_KEY_PASS" env-default:"secret"`
	Salt       string `yaml:"salt" env:"SSH_KEY_SALT" env-default:"salt"`
}

func New() *Config {
	var cfg Config
	var env string
	flag.StringVar(&env, "env", "local", "Application profile")
	flag.Parse()

	if err := cleanenv.ReadConfig(fmt.Sprintf("./config.%s.yaml", env), &cfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &cfg
}
