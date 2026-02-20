package config

import (
	"fmt"
	"time"
)

type InfraConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Minio    MinioConfig    `mapstructure:"minio"`
	Vault    VaultConfig    `mapstructure:"vault"`
}

type PostgresConfig struct {
	Platform PlatformDBConfig `mapstructure:"platform"`
	Tenant   TenantDBConfig   `mapstructure:"tenant"`
}

type PlatformDBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type TenantDBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
	Region    string `mapstructure:"region"`
}

type VaultConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Token         string `mapstructure:"token"`
	Namespace     string `mapstructure:"namespace"`
	TLSEnabled    bool   `mapstructure:"tls_enabled"`
	CACert        string `mapstructure:"ca_cert"`
	ClientCert    string `mapstructure:"client_cert"`
	ClientKey     string `mapstructure:"client_key"`
	TLSServerName string `mapstructure:"tls_server_name"`
	TLSInsecure   bool   `mapstructure:"tls_insecure"`
}

func (c *VaultConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *PlatformDBConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

func (c *TenantDBConfig) GetDSN(databaseName string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, databaseName, c.SSLMode,
	)
}

func (c *RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *MinioConfig) GetEndpoint() string {
	return c.Endpoint
}
