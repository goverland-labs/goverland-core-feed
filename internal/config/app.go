package config

type App struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	Prometheus  Prometheus
	Health      Health
	DB          DB
	Nats        Nats
	InternalAPI InternalAPI
}
