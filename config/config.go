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
	Vault      `yaml:"vault"`
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

type Vault struct {
	Address     string           `yaml:"address" env-default:"localhost"`
	Port        string           `yaml:"port" env-default:"8200"`
	Token       string           `yaml:"token" env:"VAULT_TOKEN"`
	Connections ConnectionsVault `yaml:"connection"`
}

type ConnectionsVault struct {
	Path      string `yaml:"path" env:"VAULT_PATH"`
	MountPath string `yaml:"mount-path" env:"VAULT_MOUNT_PATH"`
}

func New() *Config {
	var cfg Config
	var env string
	flag.StringVar(&env, "env", "local", "Application profile")
	flag.Parse()

	if err := cleanenv.ReadConfig(fmt.Sprintf("../config.%s.yaml", env), &cfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &cfg
}
