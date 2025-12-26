package wiring

type AppConfig struct {
	JwtKeyName string `env:"JWTKEYNAME"`
	DbConnName string
	ServerPort int16  `env:"PORT"`
	Provider   string `env:"DB_PROVIDER"`
	DbUserName string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
	Profile    string `env:"PROFILE"`
}
