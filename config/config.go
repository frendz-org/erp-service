package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Server     ServerConfig     `mapstructure:"server"`
	Infra      InfraConfig      `mapstructure:"infra"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	Email      EmailConfig      `mapstructure:"email"`
	OTP        OTPConfig        `mapstructure:"otp"`
	Password   PasswordConfig   `mapstructure:"password"`
	Masterdata MasterdataConfig `mapstructure:"masterdata"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	setDefaults()
	bindEnvVariables()

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func bindEnvVariables() {
	_ = viper.BindEnv("app.name", "APP_NAME")
	_ = viper.BindEnv("app.environment", "APP_ENV")
	_ = viper.BindEnv("app.version", "APP_VERSION")

	_ = viper.BindEnv("server.host", "SERVER_HOST")
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	_ = viper.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	_ = viper.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")
	_ = viper.BindEnv("server.cors_origins", "SERVER_CORS_ORIGINS")

	_ = viper.BindEnv("infra.postgres.platform.host", "POSTGRES_HOST")
	_ = viper.BindEnv("infra.postgres.platform.port", "POSTGRES_PORT")
	_ = viper.BindEnv("infra.postgres.platform.user", "POSTGRES_USER")
	_ = viper.BindEnv("infra.postgres.platform.password", "POSTGRES_PASSWORD")
	_ = viper.BindEnv("infra.postgres.platform.database", "POSTGRES_DB")
	_ = viper.BindEnv("infra.postgres.platform.ssl_mode", "POSTGRES_SSL_MODE")
	_ = viper.BindEnv("infra.postgres.platform.max_open_conns", "POSTGRES_MAX_OPEN_CONNS")
	_ = viper.BindEnv("infra.postgres.platform.max_idle_conns", "POSTGRES_MAX_IDLE_CONNS")
	_ = viper.BindEnv("infra.postgres.platform.conn_max_lifetime", "POSTGRES_CONN_MAX_LIFETIME")

	_ = viper.BindEnv("infra.postgres.tenant.host", "POSTGRES_TENANT_HOST", "POSTGRES_HOST")
	_ = viper.BindEnv("infra.postgres.tenant.port", "POSTGRES_TENANT_PORT", "POSTGRES_PORT")
	_ = viper.BindEnv("infra.postgres.tenant.user", "POSTGRES_TENANT_USER", "POSTGRES_USER")
	_ = viper.BindEnv("infra.postgres.tenant.password", "POSTGRES_TENANT_PASSWORD", "POSTGRES_PASSWORD")
	_ = viper.BindEnv("infra.postgres.tenant.ssl_mode", "POSTGRES_TENANT_SSL_MODE", "POSTGRES_SSL_MODE")

	_ = viper.BindEnv("infra.redis.host", "REDIS_HOST")
	_ = viper.BindEnv("infra.redis.port", "REDIS_PORT")
	_ = viper.BindEnv("infra.redis.password", "REDIS_PASSWORD")
	_ = viper.BindEnv("infra.redis.db", "REDIS_DB")
	_ = viper.BindEnv("infra.redis.pool_size", "REDIS_POOL_SIZE")
	_ = viper.BindEnv("infra.redis.min_idle_conns", "REDIS_MIN_IDLE_CONNS")
	_ = viper.BindEnv("infra.redis.conn_max_idle_time", "REDIS_CONN_MAX_IDLE_TIME")
	_ = viper.BindEnv("infra.redis.conn_max_lifetime", "REDIS_CONN_MAX_LIFETIME")
	_ = viper.BindEnv("infra.redis.read_timeout", "REDIS_READ_TIMEOUT")
	_ = viper.BindEnv("infra.redis.write_timeout", "REDIS_WRITE_TIMEOUT")

	_ = viper.BindEnv("infra.minio.endpoint", "MINIO_ENDPOINT")
	_ = viper.BindEnv("infra.minio.access_key", "MINIO_ACCESS_KEY")
	_ = viper.BindEnv("infra.minio.secret_key", "MINIO_SECRET_KEY")
	_ = viper.BindEnv("infra.minio.bucket", "MINIO_BUCKET")
	_ = viper.BindEnv("infra.minio.use_ssl", "MINIO_USE_SSL")
	_ = viper.BindEnv("infra.minio.region", "MINIO_REGION")

	_ = viper.BindEnv("infra.vault.address", "VAULT_ADDR")
	_ = viper.BindEnv("infra.vault.token", "VAULT_TOKEN")

	_ = viper.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	_ = viper.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	_ = viper.BindEnv("jwt.private_key_path", "JWT_PRIVATE_KEY_PATH")
	_ = viper.BindEnv("jwt.public_key_path", "JWT_PUBLIC_KEY_PATH")
	_ = viper.BindEnv("jwt.signing_method", "JWT_SIGNING_METHOD")
	_ = viper.BindEnv("jwt.access_expiry", "JWT_ACCESS_EXPIRY")
	_ = viper.BindEnv("jwt.refresh_expiry", "JWT_REFRESH_EXPIRY")
	_ = viper.BindEnv("jwt.issuer", "JWT_ISSUER")
	_ = viper.BindEnv("jwt.audience", "JWT_AUDIENCE")
	_ = viper.BindEnv("jwt.pin_token_expiry", "JWT_PIN_TOKEN_EXPIRY")
	_ = viper.BindEnv("jwt.registration_expiry", "JWT_REGISTRATION_EXPIRY")
	_ = viper.BindEnv("jwt.registration_secret", "JWT_REGISTRATION_SECRET")

	_ = viper.BindEnv("log.level", "LOG_LEVEL")
	_ = viper.BindEnv("log.format", "LOG_FORMAT")
	_ = viper.BindEnv("log.file_path", "LOG_FILE_PATH")
	_ = viper.BindEnv("log.max_size_mb", "LOG_MAX_SIZE_MB")
	_ = viper.BindEnv("log.max_backups", "LOG_MAX_BACKUPS")
	_ = viper.BindEnv("log.max_age_days", "LOG_MAX_AGE_DAYS")
	_ = viper.BindEnv("log.compress", "LOG_COMPRESS")
	_ = viper.BindEnv("log.retain_all", "LOG_RETAIN_ALL")

	_ = viper.BindEnv("email.provider", "EMAIL_PROVIDER")
	_ = viper.BindEnv("email.smtp_host", "EMAIL_SMTP_HOST")
	_ = viper.BindEnv("email.smtp_port", "EMAIL_SMTP_PORT")
	_ = viper.BindEnv("email.smtp_user", "EMAIL_SMTP_USER")
	_ = viper.BindEnv("email.smtp_pass", "EMAIL_SMTP_PASS")
	_ = viper.BindEnv("email.from_address", "EMAIL_FROM_ADDRESS")
	_ = viper.BindEnv("email.from_name", "EMAIL_FROM_NAME")

	_ = viper.BindEnv("masterdata.cache_ttl_categories", "MASTERDATA_CACHE_TTL_CATEGORIES")
	_ = viper.BindEnv("masterdata.cache_ttl_items", "MASTERDATA_CACHE_TTL_ITEMS")
	_ = viper.BindEnv("masterdata.cache_ttl_tree", "MASTERDATA_CACHE_TTL_TREE")
}

func setDefaults() {
	viper.SetDefault("app.name", "iam-service")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.version", "1.0.0")

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("server.idle_timeout", 120*time.Second)
	viper.SetDefault("server.cors_origins", "*")

	viper.SetDefault("infra.postgres.platform.host", "localhost")
	viper.SetDefault("infra.postgres.platform.port", 5432)
	viper.SetDefault("infra.postgres.platform.database", "iam_db")
	viper.SetDefault("infra.postgres.platform.ssl_mode", "require")
	viper.SetDefault("infra.postgres.platform.max_open_conns", 25)
	viper.SetDefault("infra.postgres.platform.max_idle_conns", 10)
	viper.SetDefault("infra.postgres.platform.conn_max_lifetime", 5*time.Minute)

	viper.SetDefault("infra.postgres.tenant.host", "localhost")
	viper.SetDefault("infra.postgres.tenant.port", 5432)
	viper.SetDefault("infra.postgres.tenant.ssl_mode", "require")
	viper.SetDefault("infra.postgres.tenant.max_open_conns", 10)
	viper.SetDefault("infra.postgres.tenant.max_idle_conns", 5)
	viper.SetDefault("infra.postgres.tenant.conn_max_lifetime", 5*time.Minute)

	viper.SetDefault("infra.redis.host", "localhost")
	viper.SetDefault("infra.redis.port", 6379)
	viper.SetDefault("infra.redis.db", 0)
	viper.SetDefault("infra.redis.pool_size", 20)
	viper.SetDefault("infra.redis.min_idle_conns", 5)
	viper.SetDefault("infra.redis.conn_max_idle_time", 5*time.Minute)
	viper.SetDefault("infra.redis.conn_max_lifetime", 30*time.Minute)
	viper.SetDefault("infra.redis.read_timeout", 3*time.Second)
	viper.SetDefault("infra.redis.write_timeout", 3*time.Second)

	viper.SetDefault("infra.minio.endpoint", "localhost:9000")
	viper.SetDefault("infra.minio.use_ssl", false)
	viper.SetDefault("infra.minio.bucket", "iam-storage")
	viper.SetDefault("infra.minio.region", "us-east-1")

	viper.SetDefault("infra.vault.address", "http://localhost:8200")

	viper.SetDefault("jwt.signing_method", "HS256")
	viper.SetDefault("jwt.access_expiry", 15*time.Minute)
	viper.SetDefault("jwt.refresh_expiry", 30*24*time.Hour)
	viper.SetDefault("jwt.issuer", "iam-service")
	viper.SetDefault("jwt.audience", []string{"backoffice", "main-app"})
	viper.SetDefault("jwt.pin_token_expiry", 10*time.Minute)
	viper.SetDefault("jwt.registration_expiry", 10*time.Minute)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.file_path", "logs/app.log")
	viper.SetDefault("log.max_size_mb", 100)
	viper.SetDefault("log.max_backups", 30)
	viper.SetDefault("log.max_age_days", 30)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("log.retain_all", false)

	viper.SetDefault("email.provider", "console")
	viper.SetDefault("email.smtp_port", 587)

	viper.SetDefault("otp.length", 6)
	viper.SetDefault("otp.expiry_minutes", 10)
	viper.SetDefault("otp.max_active_otps", 3)
	viper.SetDefault("otp.resend_cooldown", 60)
	viper.SetDefault("otp.max_resend_per_hour", 5)

	viper.SetDefault("password.min_length", 8)
	viper.SetDefault("password.require_uppercase", true)
	viper.SetDefault("password.require_lowercase", true)
	viper.SetDefault("password.require_number", true)
	viper.SetDefault("password.require_special", false)
	viper.SetDefault("password.history_count", 5)

	viper.SetDefault("masterdata.cache_ttl_categories", 24*time.Hour)
	viper.SetDefault("masterdata.cache_ttl_items", 1*time.Hour)
	viper.SetDefault("masterdata.cache_ttl_tree", 1*time.Hour)
}

func (c *Config) Validate() error {
	if c.JWT.SigningMethod == "RS256" {
		if c.JWT.PrivateKeyPath == "" {
			return fmt.Errorf("JWT_PRIVATE_KEY_PATH is required for RS256 signing")
		}
		if c.JWT.PublicKeyPath == "" {
			return fmt.Errorf("JWT_PUBLIC_KEY_PATH is required for RS256 signing")
		}
	} else if c.JWT.SigningMethod == "HS256" {
		if c.JWT.AccessSecret == "" {
			return fmt.Errorf("JWT_ACCESS_SECRET is required for HS256 signing")
		}
		if c.JWT.RefreshSecret == "" {
			return fmt.Errorf("JWT_REFRESH_SECRET is required for HS256 signing")
		}
	} else {
		return fmt.Errorf("JWT_SIGNING_METHOD must be either 'HS256' or 'RS256'")
	}

	if c.Infra.Postgres.Platform.User == "" {
		return fmt.Errorf("POSTGRES_USER is required")
	}
	if c.Infra.Postgres.Platform.Password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.App.IsDevelopment()
}

func (c *Config) IsProduction() bool {
	return c.App.IsProduction()
}
