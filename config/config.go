package config

import "github.com/spf13/viper"

type App struct {
	AppEnv  string `json:"app_env"`
	AppPort string `json:"app_port"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type Psql struct {
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBPort     string `json:"db_port"`
	DBMaxOpen  int    `json:"db_max_open"`
	DBMaxIdle  int    `json:"db_max_idle"`
}

type CloudFlareR2 struct {
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Token     string `json:"token"`
	AccountID string `json:"account_id"`
	PublicURL string `json:"public_url"`
}

type Config struct {
	App  App
	Psql Psql
	R2   CloudFlareR2
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppEnv:  viper.GetString("APP_ENV"),
			AppPort: viper.GetString("APP_PORT"),

			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),
		},
		Psql: Psql{
			DBName:     viper.GetString("DATABASE_NAME"),
			DBHost:     viper.GetString("DATABASE_HOST"),
			DBUser:     viper.GetString("DATABASE_USER"),
			DBPassword: viper.GetString("DATABASE_PASSWORD"),
			DBPort:     viper.GetString("DATABASE_PORT"),
			DBMaxOpen:  viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),
			DBMaxIdle:  viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		R2: CloudFlareR2{
			Name:      viper.GetString("CLOUDFLARE_R2_BUCKET_NAME"),
			ApiKey:    viper.GetString("CLOUDFLARE_R2_API_KEY"),
			ApiSecret: viper.GetString("CLOUDFLARE_R2_API_SECRET"),
			Token:     viper.GetString("CLOUDFLARE_R2_TOKEN"),
			AccountID: viper.GetString("CLOUDFLARE_R2_ACCOUNT_ID"),
			PublicURL: viper.GetString("CLOUDFLARE_R2_PUBLIC_URL"),
		},
	}
}
