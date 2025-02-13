package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// generall config struct
type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	TokenTTL time.Duration  `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPConfig      `yaml:"grpc"`
	Metrics  Metrics        `yaml:"metrics"`
	Jaeger   Jaeger         `yaml:"jaeger"`
}

// postgres config
type PostgresConfig struct {
	PostgresqlHost     string `yaml:"postgresql_host"`
	PostgresqlPort     string `yaml:"postgresql_port"`
	PostgresqlUser     string `yaml:"postgresql_user"`
	PostgresqlPassword string `yaml:"postgresql_password"`
	PostgresqlDbname   string `yaml:"postgresql_dbname"`
	PostgresqlSSLMode  string `yaml:"postgresql_sslmode"`
	PgDriver           string `yaml:"pg_driver"`
}

// redis config
type RedisConfig struct {
	RedisAddr      string `yaml:"redis_addr"`
	RedisPassword  string `yaml:"redis_password"`
	RedisDB        string `yaml:"redis_db"`
	RedisDefaultdb string `yaml:"redis_default_db"`
	MinIdleConns   int    `yaml:"redis_min_idle_conns"`
	PoolSize       int    `yaml:"redis_pool_size"`
	PoolTimeout    int    `yaml:"redis_pool_timeout"`
	Password       string `yaml:"rd_password"`
	DB             int    `yaml:"rd_db"`
}

// grpc config
type GRPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Metrics struct {
	Url         string `yaml:"prom_url"`
	ServiceName string `yaml:"prom_service_name"`
}

type Jaeger struct {
	Host        string `yaml:"jaeger_host"`
	ServiceName string `yaml:"jaeger_service_name"`
	LogSpans    bool   `yaml:"jaeger_log_spans"`
}

var (
	ErrInvalidOsEnvironmentspssw = errors.New("op cannot find variables. Password")
	ErrInvalidOsEnvironmentsuser = errors.New("op cannot find variables. User")
)

func MustLoad() *Config {
	path := FetchConfigFlag()
	if path == "" {
		panic("config path is empty" + path)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is not exist" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config" + err.Error())
	}

	if err := FetchVariables(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func FetchConfigFlag() string {
	var res string

	flag.StringVar(&res, "config", "", "path to the file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}

func FetchVariables(cfg *Config) error {

	cfg.Postgres.PostgresqlPassword = os.Getenv("POSTGRES_PASSWORD")
	if cfg.Postgres.PostgresqlPassword == "" {
		return ErrInvalidOsEnvironmentspssw
	}

	cfg.Postgres.PostgresqlUser = os.Getenv("POSTGRES_USER")
	if cfg.Postgres.PostgresqlPassword == "" {
		return ErrInvalidOsEnvironmentsuser
	}

	return nil
}
