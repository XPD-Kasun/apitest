package wiring

type AppConfig struct {
	EnvKeyName string
	DbConnName string
	ServerPort int16
	Provider   string
}
