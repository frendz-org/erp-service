package config

import "time"

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	CORSOrigins  string        `mapstructure:"cors_origins"`
}

type JWTConfig struct {
	AccessSecret       string `mapstructure:"access_secret"`
	RefreshSecret      string `mapstructure:"refresh_secret"`
	RegistrationSecret string `mapstructure:"registration_secret"`

	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKeyPath  string `mapstructure:"public_key_path"`
	SigningMethod  string `mapstructure:"signing_method"`

	AccessExpiry       time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry      time.Duration `mapstructure:"refresh_expiry"`
	Issuer             string        `mapstructure:"issuer"`
	Audience           []string      `mapstructure:"audience"`
	PINTokenExpiry     time.Duration `mapstructure:"pin_token_expiry"`
	RegistrationExpiry time.Duration `mapstructure:"registration_expiry"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`

	AuditEnabled bool `mapstructure:"audit_enabled"`

	FilePath   string `mapstructure:"file_path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
	Compress   bool   `mapstructure:"compress"`
	RetainAll  bool   `mapstructure:"retain_all"`
}

type EmailConfig struct {
	Provider    string `mapstructure:"provider"`
	SMTPHost    string `mapstructure:"smtp_host"`
	SMTPPort    int    `mapstructure:"smtp_port"`
	SMTPUser    string `mapstructure:"smtp_user"`
	SMTPPass    string `mapstructure:"smtp_pass"`
	FromAddress string `mapstructure:"from_address"`
	FromName    string `mapstructure:"from_name"`
}

type OTPConfig struct {
	Length           int `mapstructure:"length"`
	ExpiryMinutes    int `mapstructure:"expiry_minutes"`
	MaxActiveOTPs    int `mapstructure:"max_active_otps"`
	ResendCooldown   int `mapstructure:"resend_cooldown"`
	MaxResendPerHour int `mapstructure:"max_resend_per_hour"`
}

type PasswordConfig struct {
	MinLength        int  `mapstructure:"min_length"`
	RequireUppercase bool `mapstructure:"require_uppercase"`
	RequireLowercase bool `mapstructure:"require_lowercase"`
	RequireNumber    bool `mapstructure:"require_number"`
	RequireSpecial   bool `mapstructure:"require_special"`
	HistoryCount     int  `mapstructure:"history_count"`
}

type MasterdataConfig struct {
	CacheTTLCategories time.Duration `mapstructure:"cache_ttl_categories"`
	CacheTTLItems      time.Duration `mapstructure:"cache_ttl_items"`
	CacheTTLTree       time.Duration `mapstructure:"cache_ttl_tree"`
}

type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

func (c *GoogleOAuthConfig) IsEnabled() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

func (c *AppConfig) IsStaging() bool {
	return c.Environment == "staging"
}
