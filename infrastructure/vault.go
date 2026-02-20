package infrastructure

import (
	"fmt"
	"log"

	"iam-service/config"

	"github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Address   string
	Token     string
	Namespace string
	TLSConfig *api.TLSConfig
}

func NewVault(cfg config.VaultConfig) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.GetAddress()

	if cfg.TLSEnabled {
		tlsConfig := &api.TLSConfig{
			CACert:        cfg.CACert,
			ClientCert:    cfg.ClientCert,
			ClientKey:     cfg.ClientKey,
			TLSServerName: cfg.TLSServerName,
			Insecure:      cfg.TLSInsecure,
		}
		if err := vaultConfig.ConfigureTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
	}

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	if cfg.Token != "" {
		client.SetToken(cfg.Token)
	}

	if cfg.Namespace != "" {
		client.SetNamespace(cfg.Namespace)
	}

	if cfg.Token != "" {
		if _, err := client.Auth().Token().LookupSelf(); err != nil {
			return nil, fmt.Errorf("failed to verify vault connection: %w", err)
		}
		log.Printf("Connected to Vault at %s\n", cfg.GetAddress())
	} else {
		log.Printf("Vault client created for %s (no token set)\n", cfg.GetAddress())
	}

	return client, nil
}

func NewVaultWithAppRole(cfg config.VaultConfig, roleID, secretID string) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.GetAddress()

	if cfg.TLSEnabled {
		tlsConfig := &api.TLSConfig{
			CACert:        cfg.CACert,
			ClientCert:    cfg.ClientCert,
			ClientKey:     cfg.ClientKey,
			TLSServerName: cfg.TLSServerName,
			Insecure:      cfg.TLSInsecure,
		}
		if err := vaultConfig.ConfigureTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
	}

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	if cfg.Namespace != "" {
		client.SetNamespace(cfg.Namespace)
	}

	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	secret, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with approle: %w", err)
	}

	if secret == nil || secret.Auth == nil {
		return nil, fmt.Errorf("empty response from vault approle login")
	}

	client.SetToken(secret.Auth.ClientToken)

	log.Printf("Connected to Vault at %s using AppRole\n", cfg.GetAddress())

	return client, nil
}

func NewVaultWithKubernetes(cfg config.VaultConfig, role, jwt string) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.GetAddress()

	if cfg.TLSEnabled {
		tlsConfig := &api.TLSConfig{
			CACert:        cfg.CACert,
			ClientCert:    cfg.ClientCert,
			ClientKey:     cfg.ClientKey,
			TLSServerName: cfg.TLSServerName,
			Insecure:      cfg.TLSInsecure,
		}
		if err := vaultConfig.ConfigureTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
	}

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	if cfg.Namespace != "" {
		client.SetNamespace(cfg.Namespace)
	}

	data := map[string]interface{}{
		"role": role,
		"jwt":  jwt,
	}

	secret, err := client.Logical().Write("auth/kubernetes/login", data)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with kubernetes: %w", err)
	}

	if secret == nil || secret.Auth == nil {
		return nil, fmt.Errorf("empty response from vault kubernetes login")
	}

	client.SetToken(secret.Auth.ClientToken)

	log.Printf("Connected to Vault at %s using Kubernetes auth\n", cfg.GetAddress())

	return client, nil
}
