package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
	runtime "runtime"
	"time"
)

type Config struct {
	App        `yaml:"app"`
	HTTPServer `yaml:"http-server"`
	DB         `yaml:"database"`
	Vault      `yaml:"vault"`
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

type Vault struct {
	Address     string           `yaml:"address" env-default:"localhost"`
	Port        string           `yaml:"port" env-default:"8200"`
	Token       string           `yaml:"token" env:"VAULT_TOKEN"`
	MountPath   string           `yaml:"mount-path" env:"VAULT_MOUNT_PATH"`
	Connections ConnectionsVault `yaml:"connection"`
}

type ConnectionsVault struct {
	Path string `yaml:"path" env:"VAULT_CONNECTIONS_PATH"`
}

func New() *Config {
	var env string
	flag.StringVar(&env, "env", "local", "Application profile")
	flag.Parse()

	return NewWithEnvironment(env)
}

func NewWithEnvironment(env string) *Config {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get caller info")
	}

	projectRoot := filepath.Dir(filepath.Dir(filename))
	configPath := filepath.Join(projectRoot, fmt.Sprintf("config.%s.yaml", env))

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &cfg
}
