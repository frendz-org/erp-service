package hashivault

import (
	"context"
	"erp-service/pkg/errors"
	"fmt"
	"time"
)

type DatabaseCredentials struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
}

type DatabaseRole struct {
	DBName               string   `json:"db_name"`
	CreationStatements   []string `json:"creation_statements"`
	RevocationStatements []string `json:"revocation_statements"`
	RollbackStatements   []string `json:"rollback_statements"`
	RenewStatements      []string `json:"renew_statements"`
	DefaultTTL           int      `json:"default_ttl"`
	MaxTTL               int      `json:"max_ttl"`
}

func (v *SecureVault) GenerateDatabaseCredentials(ctx context.Context, role string) (*DatabaseCredentials, error) {
	path := fmt.Sprintf("database/creds/%s", role)
	secret, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate database credentials").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("vault returned empty response")
	}

	creds := &DatabaseCredentials{
		LeaseID:       secret.LeaseID,
		LeaseDuration: secret.LeaseDuration,
		Renewable:     secret.Renewable,
	}

	if v, ok := secret.Data["username"].(string); ok {
		creds.Username = v
	}
	if v, ok := secret.Data["password"].(string); ok {
		creds.Password = v
	}

	return creds, nil
}

func (v *SecureVault) GenerateStaticDatabaseCredentials(ctx context.Context, role string) (*DatabaseCredentials, error) {
	path := fmt.Sprintf("database/static-creds/%s", role)
	secret, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, errors.ErrInternal("failed to get static database credentials").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("vault returned empty response")
	}

	creds := &DatabaseCredentials{}

	if v, ok := secret.Data["username"].(string); ok {
		creds.Username = v
	}
	if v, ok := secret.Data["password"].(string); ok {
		creds.Password = v
	}
	if v, ok := secret.Data["last_vault_rotation"].(string); ok {

		_ = v
	}
	if v, ok := secret.Data["rotation_period"].(int); ok {
		_ = v
	}

	return creds, nil
}

func (v *SecureVault) RotateStaticDatabaseCredentials(ctx context.Context, role string) error {
	path := fmt.Sprintf("database/rotate-role/%s", role)
	_, err := v.client.Logical().WriteWithContext(ctx, path, nil)
	if err != nil {
		return errors.ErrInternal("failed to rotate static credentials").WithError(err)
	}
	return nil
}

func (v *SecureVault) CreateDatabaseRole(ctx context.Context, name string, role *DatabaseRole) error {
	data := map[string]interface{}{
		"db_name":             role.DBName,
		"creation_statements": role.CreationStatements,
		"default_ttl":         role.DefaultTTL,
		"max_ttl":             role.MaxTTL,
	}

	if len(role.RevocationStatements) > 0 {
		data["revocation_statements"] = role.RevocationStatements
	}
	if len(role.RollbackStatements) > 0 {
		data["rollback_statements"] = role.RollbackStatements
	}
	if len(role.RenewStatements) > 0 {
		data["renew_statements"] = role.RenewStatements
	}

	path := fmt.Sprintf("database/roles/%s", name)
	_, err := v.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return errors.ErrInternal("failed to create database role").WithError(err)
	}

	return nil
}

func (v *SecureVault) ReadDatabaseRole(ctx context.Context, name string) (*DatabaseRole, error) {
	path := fmt.Sprintf("database/roles/%s", name)
	secret, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, errors.ErrInternal("failed to read database role").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("database role not found")
	}

	role := &DatabaseRole{}

	if v, ok := secret.Data["db_name"].(string); ok {
		role.DBName = v
	}
	if v, ok := secret.Data["default_ttl"].(int); ok {
		role.DefaultTTL = v
	} else if v, ok := secret.Data["default_ttl"].(float64); ok {
		role.DefaultTTL = int(v)
	}
	if v, ok := secret.Data["max_ttl"].(int); ok {
		role.MaxTTL = v
	} else if v, ok := secret.Data["max_ttl"].(float64); ok {
		role.MaxTTL = int(v)
	}

	if v, ok := secret.Data["creation_statements"].([]interface{}); ok {
		statements := make([]string, 0, len(v))
		for _, s := range v {
			if str, ok := s.(string); ok {
				statements = append(statements, str)
			}
		}
		role.CreationStatements = statements
	}

	if v, ok := secret.Data["revocation_statements"].([]interface{}); ok {
		statements := make([]string, 0, len(v))
		for _, s := range v {
			if str, ok := s.(string); ok {
				statements = append(statements, str)
			}
		}
		role.RevocationStatements = statements
	}

	return role, nil
}

func (v *SecureVault) DeleteDatabaseRole(ctx context.Context, name string) error {
	path := fmt.Sprintf("database/roles/%s", name)
	_, err := v.client.Logical().DeleteWithContext(ctx, path)
	if err != nil {
		return errors.ErrInternal("failed to delete database role").WithError(err)
	}
	return nil
}

func (v *SecureVault) ListDatabaseRoles(ctx context.Context) ([]string, error) {
	secret, err := v.client.Logical().ListWithContext(ctx, "database/roles")
	if err != nil {
		return nil, errors.ErrInternal("failed to list database roles").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	keysInterface, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	keysSlice, ok := keysInterface.([]interface{})
	if !ok {
		return nil, errors.ErrInternal("unexpected keys format")
	}

	keys := make([]string, 0, len(keysSlice))
	for _, k := range keysSlice {
		if keyStr, ok := k.(string); ok {
			keys = append(keys, keyStr)
		}
	}

	return keys, nil
}

func (v *SecureVault) ConfigureDatabaseConnection(ctx context.Context, name, plugin string, config map[string]interface{}) error {
	data := map[string]interface{}{
		"plugin_name": plugin,
	}

	for k, v := range config {
		data[k] = v
	}

	path := fmt.Sprintf("database/config/%s", name)
	_, err := v.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return errors.ErrInternal("failed to configure database connection").WithError(err)
	}

	return nil
}

func (v *SecureVault) ReadDatabaseConnection(ctx context.Context, name string) (map[string]interface{}, error) {
	path := fmt.Sprintf("database/config/%s", name)
	secret, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, errors.ErrInternal("failed to read database connection").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("database connection not found")
	}

	return secret.Data, nil
}

func (v *SecureVault) DeleteDatabaseConnection(ctx context.Context, name string) error {
	path := fmt.Sprintf("database/config/%s", name)
	_, err := v.client.Logical().DeleteWithContext(ctx, path)
	if err != nil {
		return errors.ErrInternal("failed to delete database connection").WithError(err)
	}
	return nil
}

func (v *SecureVault) ResetDatabaseConnection(ctx context.Context, name string) error {
	path := fmt.Sprintf("database/reset/%s", name)
	_, err := v.client.Logical().WriteWithContext(ctx, path, nil)
	if err != nil {
		return errors.ErrInternal("failed to reset database connection").WithError(err)
	}
	return nil
}

func (v *SecureVault) RotateDatabaseRootCredentials(ctx context.Context, name string) error {
	path := fmt.Sprintf("database/rotate-root/%s", name)
	_, err := v.client.Logical().WriteWithContext(ctx, path, nil)
	if err != nil {
		return errors.ErrInternal("failed to rotate root credentials").WithError(err)
	}
	return nil
}

func (v *SecureVault) ConfigurePostgreSQLConnection(ctx context.Context, name, connectionURL string, maxOpenConns, maxIdleConns int, maxConnLifetime time.Duration) error {
	config := map[string]interface{}{
		"connection_url": connectionURL,
	}

	if maxOpenConns > 0 {
		config["max_open_connections"] = maxOpenConns
	}
	if maxIdleConns > 0 {
		config["max_idle_connections"] = maxIdleConns
	}
	if maxConnLifetime > 0 {
		config["max_connection_lifetime"] = maxConnLifetime.String()
	}

	return v.ConfigureDatabaseConnection(ctx, name, "postgresql-database-plugin", config)
}

func (v *SecureVault) ConfigureMySQLConnection(ctx context.Context, name, connectionURL string, maxOpenConns, maxIdleConns int, maxConnLifetime time.Duration) error {
	config := map[string]interface{}{
		"connection_url": connectionURL,
	}

	if maxOpenConns > 0 {
		config["max_open_connections"] = maxOpenConns
	}
	if maxIdleConns > 0 {
		config["max_idle_connections"] = maxIdleConns
	}
	if maxConnLifetime > 0 {
		config["max_connection_lifetime"] = maxConnLifetime.String()
	}

	return v.ConfigureDatabaseConnection(ctx, name, "mysql-database-plugin", config)
}
