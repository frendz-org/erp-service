package hashivault

import (
	"github.com/hashicorp/vault/api"
)

type SecureVault struct {
	client *api.Client
}

func NewSecureVault(client *api.Client) *SecureVault {
	return &SecureVault{
		client: client,
	}
}
