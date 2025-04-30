package config

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config структура, содержащая всю конфигурацию
type Config struct {
	Env       string `env:"ENV" env-required:"true"`
	WebSocket WebSocketConfig
	Http      HttpConfig
	Talk      TalkConfig
	Sso       SSOConfig
	Postgres  PostgresConfig
}

// WebSocketConfig структура, содержащая настройки для WebSocket сервера
type WebSocketConfig struct {
	Port string `env:"WEBSOCKET_PORT" env-required:"true"`
}

// WebSocketConfig структура, содержащая настройки для HTTP сервера
type HttpConfig struct {
	Host string `env:"HTTP_HOST" env-required:"true"`
	Port string `env:"HTTP_PORT" env-required:"true"`
}

// TalkConfig структура, содержащая вспомогательные данные для работы текущего сервиса
type TalkConfig struct {
	HttpClientUrl string `env:"TALK_HTTP_CLIENT_URL" env-required:"true"`
}

// SSOConfig структура, содержащая настройки для подключения к сервису SSO
type SSOConfig struct {
	HttpServerUrl string `env:"SSO_HTTP_SERVER_URL" env-required:"true"`
	HttpClientUrl string `env:"SSO_HTTP_CLIENT_URL" env-required:"true"`
	GrpcServerUrl string `env:"SSO_GRPC_SERVER_URL" env-required:"true"`
}

// PostgresConfig структура, содержащая настройки для подключения к Postgresql
type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Database string `env:"POSTGRES_DATABASE" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
}

var config Config

// Get возвращает копию текущей конфигурации
func Get() Config {
	return config
}

// MustLoad загружает конфигурацию из файла и возвращает её
func MustLoad() Config {
	// Чтение переменных из окружения
	err := cleanenv.ReadEnv(&config)
	if err == nil {
		return config
	} else {
		fmt.Println("error on read raw env: " + err.Error())
	}

	// Чтение переменных из конфигурационного файла
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty. you need to specify --config=<file_path> or environment CONFIG_PATH")
	}

	fullConfigPath := getAbsoluteConfigPath(configPath)

	if _, err := os.Stat(fullConfigPath); os.IsNotExist(err) {
		panic("config file does not exist: " + fullConfigPath)
	}

	if err := cleanenv.ReadConfig(fullConfigPath, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return config
}

// fetchConfigPath возвращает путь к конфигурационному файлу из аргументов командной строки или переменной окружения
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func getAbsoluteConfigPath(configPath string) string {
	wdFromEnv := os.Getenv("WORKDIR_PATH")
	wdFromOs, _ := os.Getwd()

	if wdFromEnv != "" {
		return path.Join(wdFromEnv, configPath)
	}
	return path.Join(wdFromOs, configPath)
}
