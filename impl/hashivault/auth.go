package hashivault

import (
	"context"
	"fmt"
	"iam-service/pkg/errors"
	"time"
)

type TokenInfo struct {
	Accessor       string            `json:"accessor"`
	CreationTime   int64             `json:"creation_time"`
	CreationTTL    int64             `json:"creation_ttl"`
	DisplayName    string            `json:"display_name"`
	EntityID       string            `json:"entity_id"`
	ExpireTime     time.Time         `json:"expire_time"`
	ExplicitMaxTTL int64             `json:"explicit_max_ttl"`
	ID             string            `json:"id"`
	Meta           map[string]string `json:"meta"`
	NumUses        int               `json:"num_uses"`
	Orphan         bool              `json:"orphan"`
	Path           string            `json:"path"`
	Policies       []string          `json:"policies"`
	Renewable      bool              `json:"renewable"`
	TTL            int64             `json:"ttl"`
	Type           string            `json:"type"`
}

func (v *SecureVault) CreateToken(ctx context.Context, options map[string]interface{}) (string, error) {
	secret, err := v.client.Logical().WriteWithContext(ctx, "auth/token/create", options)
	if err != nil {
		return "", errors.ErrInternal("failed to create token").WithError(err)
	}

	if secret == nil || secret.Auth == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	return secret.Auth.ClientToken, nil
}

func (v *SecureVault) CreateOrphanToken(ctx context.Context, options map[string]interface{}) (string, error) {
	if options == nil {
		options = make(map[string]interface{})
	}
	options["no_parent"] = true

	secret, err := v.client.Logical().WriteWithContext(ctx, "auth/token/create-orphan", options)
	if err != nil {
		return "", errors.ErrInternal("failed to create orphan token").WithError(err)
	}

	if secret == nil || secret.Auth == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	return secret.Auth.ClientToken, nil
}

func (v *SecureVault) LookupToken(ctx context.Context, token string) (*TokenInfo, error) {
	data := map[string]interface{}{
		"token": token,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "auth/token/lookup", data)
	if err != nil {
		return nil, errors.ErrInternal("failed to lookup token").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("token not found")
	}

	return parseTokenInfo(secret.Data)
}

func (v *SecureVault) LookupSelfToken(ctx context.Context) (*TokenInfo, error) {
	secret, err := v.client.Logical().ReadWithContext(ctx, "auth/token/lookup-self")
	if err != nil {
		return nil, errors.ErrInternal("failed to lookup self token").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("token not found")
	}

	return parseTokenInfo(secret.Data)
}

func (v *SecureVault) RenewToken(ctx context.Context, token string, increment int) (*TokenInfo, error) {
	data := map[string]interface{}{
		"token": token,
	}

	if increment > 0 {
		data["increment"] = increment
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "auth/token/renew", data)
	if err != nil {
		return nil, errors.ErrInternal("failed to renew token").WithError(err)
	}

	if secret == nil || secret.Auth == nil {
		return nil, errors.ErrNotFound("token not found or not renewable")
	}

	return &TokenInfo{
		ID:        secret.Auth.ClientToken,
		Policies:  secret.Auth.Policies,
		Renewable: secret.Auth.Renewable,
		TTL:       int64(secret.Auth.LeaseDuration),
	}, nil
}

func (v *SecureVault) RenewSelfToken(ctx context.Context, increment int) error {
	data := map[string]interface{}{}
	if increment > 0 {
		data["increment"] = increment
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "auth/token/renew-self", data)
	if err != nil {
		return errors.ErrInternal("failed to renew self token").WithError(err)
	}

	return nil
}

func (v *SecureVault) RevokeToken(ctx context.Context, token string) error {
	data := map[string]interface{}{
		"token": token,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "auth/token/revoke", data)
	if err != nil {
		return errors.ErrInternal("failed to revoke token").WithError(err)
	}

	return nil
}

func (v *SecureVault) RevokeSelfToken(ctx context.Context) error {
	_, err := v.client.Logical().WriteWithContext(ctx, "auth/token/revoke-self", nil)
	if err != nil {
		return errors.ErrInternal("failed to revoke self token").WithError(err)
	}

	return nil
}

func (v *SecureVault) RevokeOrphanToken(ctx context.Context, token string) error {
	data := map[string]interface{}{
		"token": token,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "auth/token/revoke-orphan", data)
	if err != nil {
		return errors.ErrInternal("failed to revoke orphan token").WithError(err)
	}

	return nil
}

func (v *SecureVault) RevokeTokenByAccessor(ctx context.Context, accessor string) error {
	data := map[string]interface{}{
		"accessor": accessor,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "auth/token/revoke-accessor", data)
	if err != nil {
		return errors.ErrInternal("failed to revoke token by accessor").WithError(err)
	}

	return nil
}

func (v *SecureVault) LoginWithAppRole(ctx context.Context, roleID, secretID string) (string, error) {
	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "auth/approle/login", data)
	if err != nil {
		return "", errors.ErrInternal("failed to login with approle").WithError(err)
	}

	if secret == nil || secret.Auth == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	v.client.SetToken(secret.Auth.ClientToken)

	return secret.Auth.ClientToken, nil
}

func (v *SecureVault) LoginWithUserpass(ctx context.Context, username, password string) (string, error) {
	data := map[string]interface{}{
		"password": password,
	}

	path := fmt.Sprintf("auth/userpass/login/%s", username)
	secret, err := v.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return "", errors.ErrInternal("failed to login with userpass").WithError(err)
	}

	if secret == nil || secret.Auth == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	v.client.SetToken(secret.Auth.ClientToken)

	return secret.Auth.ClientToken, nil
}

func (v *SecureVault) LoginWithKubernetes(ctx context.Context, role, jwt string) (string, error) {
	data := map[string]interface{}{
		"role": role,
		"jwt":  jwt,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "auth/kubernetes/login", data)
	if err != nil {
		return "", errors.ErrInternal("failed to login with kubernetes").WithError(err)
	}

	if secret == nil || secret.Auth == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	v.client.SetToken(secret.Auth.ClientToken)

	return secret.Auth.ClientToken, nil
}

func (v *SecureVault) SetToken(token string) {
	v.client.SetToken(token)
}

func (v *SecureVault) GetToken() string {
	return v.client.Token()
}

func (v *SecureVault) ClearToken() {
	v.client.ClearToken()
}

func parseTokenInfo(data map[string]interface{}) (*TokenInfo, error) {
	info := &TokenInfo{}

	if v, ok := data["accessor"].(string); ok {
		info.Accessor = v
	}
	if v, ok := data["creation_time"].(int64); ok {
		info.CreationTime = v
	} else if v, ok := data["creation_time"].(float64); ok {
		info.CreationTime = int64(v)
	}
	if v, ok := data["creation_ttl"].(int64); ok {
		info.CreationTTL = v
	} else if v, ok := data["creation_ttl"].(float64); ok {
		info.CreationTTL = int64(v)
	}
	if v, ok := data["display_name"].(string); ok {
		info.DisplayName = v
	}
	if v, ok := data["entity_id"].(string); ok {
		info.EntityID = v
	}
	if v, ok := data["explicit_max_ttl"].(int64); ok {
		info.ExplicitMaxTTL = v
	} else if v, ok := data["explicit_max_ttl"].(float64); ok {
		info.ExplicitMaxTTL = int64(v)
	}
	if v, ok := data["id"].(string); ok {
		info.ID = v
	}
	if v, ok := data["num_uses"].(int); ok {
		info.NumUses = v
	} else if v, ok := data["num_uses"].(float64); ok {
		info.NumUses = int(v)
	}
	if v, ok := data["orphan"].(bool); ok {
		info.Orphan = v
	}
	if v, ok := data["path"].(string); ok {
		info.Path = v
	}
	if v, ok := data["renewable"].(bool); ok {
		info.Renewable = v
	}
	if v, ok := data["ttl"].(int64); ok {
		info.TTL = v
	} else if v, ok := data["ttl"].(float64); ok {
		info.TTL = int64(v)
	}
	if v, ok := data["type"].(string); ok {
		info.Type = v
	}

	if policiesRaw, ok := data["policies"]; ok {
		if policiesSlice, ok := policiesRaw.([]interface{}); ok {
			policies := make([]string, 0, len(policiesSlice))
			for _, p := range policiesSlice {
				if pStr, ok := p.(string); ok {
					policies = append(policies, pStr)
				}
			}
			info.Policies = policies
		}
	}

	if metaRaw, ok := data["meta"]; ok {
		if metaMap, ok := metaRaw.(map[string]interface{}); ok {
			meta := make(map[string]string)
			for k, v := range metaMap {
				if vStr, ok := v.(string); ok {
					meta[k] = vStr
				}
			}
			info.Meta = meta
		}
	}

	return info, nil
}
