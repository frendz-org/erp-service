package hashivault

import (
	"context"
	"fmt"
	"iam-service/pkg/errors"
	"time"
)

type LeaseInfo struct {
	LeaseID       string    `json:"lease_id"`
	LeaseDuration int       `json:"lease_duration"`
	Renewable     bool      `json:"renewable"`
	IssueTime     time.Time `json:"issue_time"`
	ExpireTime    time.Time `json:"expire_time"`
}

func (v *SecureVault) RenewLease(ctx context.Context, leaseID string, increment int) (*LeaseInfo, error) {
	data := map[string]interface{}{
		"lease_id": leaseID,
	}

	if increment > 0 {
		data["increment"] = increment
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "sys/leases/renew", data)
	if err != nil {
		return nil, errors.ErrInternal("failed to renew lease").WithError(err)
	}

	if secret == nil {
		return nil, errors.ErrNotFound("vault returned empty response")
	}

	return &LeaseInfo{
		LeaseID:       secret.LeaseID,
		LeaseDuration: secret.LeaseDuration,
		Renewable:     secret.Renewable,
	}, nil
}

func (v *SecureVault) RevokeLease(ctx context.Context, leaseID string) error {
	data := map[string]interface{}{
		"lease_id": leaseID,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "sys/leases/revoke", data)
	if err != nil {
		return errors.ErrInternal("failed to revoke lease").WithError(err)
	}

	return nil
}

func (v *SecureVault) RevokeLeaseWithPrefix(ctx context.Context, prefix string) error {
	data := map[string]interface{}{
		"prefix": prefix,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "sys/leases/revoke-prefix", data)
	if err != nil {
		return errors.ErrInternal("failed to revoke leases with prefix").WithError(err)
	}

	return nil
}

func (v *SecureVault) ForceRevokeLeaseWithPrefix(ctx context.Context, prefix string) error {
	data := map[string]interface{}{
		"prefix": prefix,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, "sys/leases/revoke-force", data)
	if err != nil {
		return errors.ErrInternal("failed to force revoke leases with prefix").WithError(err)
	}

	return nil
}

func (v *SecureVault) LookupLease(ctx context.Context, leaseID string) (*LeaseInfo, error) {
	data := map[string]interface{}{
		"lease_id": leaseID,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "sys/leases/lookup", data)
	if err != nil {
		return nil, errors.ErrInternal("failed to lookup lease").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("vault returned empty response")
	}

	info := &LeaseInfo{}

	if v, ok := secret.Data["id"].(string); ok {
		info.LeaseID = v
	}
	if v, ok := secret.Data["ttl"].(int); ok {
		info.LeaseDuration = v
	} else if v, ok := secret.Data["ttl"].(float64); ok {
		info.LeaseDuration = int(v)
	}
	if v, ok := secret.Data["renewable"].(bool); ok {
		info.Renewable = v
	}
	if v, ok := secret.Data["issue_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			info.IssueTime = t
		}
	}
	if v, ok := secret.Data["expire_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			info.ExpireTime = t
		}
	}

	return info, nil
}

func (v *SecureVault) ListLeases(ctx context.Context, prefix string) ([]string, error) {
	path := fmt.Sprintf("sys/leases/lookup/%s", prefix)
	secret, err := v.client.Logical().ListWithContext(ctx, path)
	if err != nil {
		return nil, errors.ErrInternal("failed to list leases").WithError(err)
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

func (v *SecureVault) TidyLeases(ctx context.Context) error {
	_, err := v.client.Logical().WriteWithContext(ctx, "sys/leases/tidy", nil)
	if err != nil {
		return errors.ErrInternal("failed to tidy leases").WithError(err)
	}
	return nil
}
