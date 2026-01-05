package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SettingProvider interface {
	GetEnvironment() string
	GetServer() Server
	GetDatabase() Database
	GetBroker() Broker
	GetPerspectiveClient() PerspectiveAPIClient
	IsProd() bool
	IsTest() bool
	IsLocal() bool
}

type Database struct {
	Host    string
	Port    string
	User    string
	Pwd     string
	Name    string
	SSLMode string
}

type Server struct {
	Port string
}

type Broker struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Setting struct {
	Environment          string
	Server               Server
	Database             Database
	Broker               Broker
	PerspectiveAPIClient PerspectiveAPIClient
}

type PerspectiveAPIClient struct {
	URL   string
	TOKEN string
}

var AppSetting SettingProvider

func GetSetting() SettingProvider {
	AppSetting = loadSetting()

	return AppSetting
}

func loadSetting() *Setting {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Arquivo %s não encontrado, tentando .env padrão", envFile)
		_ = godotenv.Load()
		// se não achar nada, usa do sistema mesmo
	}

	return &Setting{
		Environment: env,
		Server: Server{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
		Database: Database{
			Host:    getEnvOrDefault("DB_HOST", "localhost"),
			Port:    getEnvOrDefault("DB_PORT", "5432"),
			User:    getEnvOrDefault("DB_USER", "postgres"),
			Pwd:     getEnvOrDefault("DB_PASSWORD", "1234"),
			Name:    getEnvOrDefault("DB_NAME", "go-api"),
			SSLMode: getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
		Broker: Broker{
			Host:     getEnvOrDefault("AMQP_HOST", "localhost"),
			Port:     getEnvOrDefault("AMQP_PORT", "5672"),
			User:     getEnvOrDefault("AMQP_USER", "guest"),
			Password: getEnvOrDefault("AMQP_PASSWORD", "guest"),
		},
		PerspectiveAPIClient: PerspectiveAPIClient{
			URL:   getEnvOrDefault("GOOGLE_PERSPECTIVE_API_BASE_URL", "https://commentanalyzer.googleapis.com/v1alpha1/comments:analyze"),
			TOKEN: getEnvOrDefault("GOOGLE_PERSPECTIVE_API_TOKEN", ""),
		},
	}
}

func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func (s *Setting) GetEnvironment() string                     { return s.Environment }
func (s *Setting) GetServer() Server                          { return s.Server }
func (s *Setting) GetDatabase() Database                      { return s.Database }
func (s *Setting) GetBroker() Broker                          { return s.Broker }
func (s *Setting) GetPerspectiveClient() PerspectiveAPIClient { return s.PerspectiveAPIClient }
func (s *Setting) IsProd() bool                               { return s.Environment == "production" }
func (s *Setting) IsTest() bool                               { return s.Environment == "test" }
func (s *Setting) IsLocal() bool                              { return s.Environment == "local" || s.Environment == "development" }
