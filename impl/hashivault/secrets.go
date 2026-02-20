package hashivault

import (
	"context"
	"encoding/json"
	"fmt"
	"iam-service/pkg/errors"
)

type SecretData struct {
	Data     map[string]interface{} `json:"data"`
	Metadata SecretMetadata         `json:"metadata"`
}

type SecretMetadata struct {
	Version        int               `json:"version"`
	CreatedTime    string            `json:"created_time"`
	DeletionTime   string            `json:"deletion_time,omitempty"`
	Destroyed      bool              `json:"destroyed"`
	CustomMetadata map[string]string `json:"custom_metadata,omitempty"`
}

func (v *SecureVault) WriteSecret(ctx context.Context, path string, data map[string]interface{}) (*SecretMetadata, error) {
	secretData := map[string]interface{}{
		"data": data,
	}

	secret, err := v.client.KVv2(getMountPath(path)).Put(ctx, getSecretPath(path), secretData)
	if err != nil {
		return nil, errors.ErrInternal("failed to write secret").WithError(err)
	}

	metadata := &SecretMetadata{
		Version:     secret.VersionMetadata.Version,
		CreatedTime: secret.VersionMetadata.CreatedTime.String(),
		Destroyed:   secret.VersionMetadata.Destroyed,
	}

	if secret.CustomMetadata != nil {
		metadata.CustomMetadata = convertToStringMap(secret.CustomMetadata)
	}

	return metadata, nil
}

func (v *SecureVault) ReadSecret(ctx context.Context, path string) (map[string]interface{}, error) {
	secret, err := v.client.KVv2(getMountPath(path)).Get(ctx, getSecretPath(path))
	if err != nil {
		return nil, errors.ErrInternal("failed to read secret").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("secret not found")
	}

	return secret.Data, nil
}

func (v *SecureVault) ReadSecretVersion(ctx context.Context, path string, version int) (map[string]interface{}, error) {
	secret, err := v.client.KVv2(getMountPath(path)).GetVersion(ctx, getSecretPath(path), version)
	if err != nil {
		return nil, errors.ErrInternal("failed to read secret version").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("secret not found")
	}

	return secret.Data, nil
}

func (v *SecureVault) ReadSecretWithMetadata(ctx context.Context, path string) (*SecretData, error) {
	secret, err := v.client.KVv2(getMountPath(path)).Get(ctx, getSecretPath(path))
	if err != nil {
		return nil, errors.ErrInternal("failed to read secret").WithError(err)
	}

	if secret == nil {
		return nil, errors.ErrNotFound("secret not found")
	}

	metadata := SecretMetadata{
		Version:     secret.VersionMetadata.Version,
		CreatedTime: secret.VersionMetadata.CreatedTime.String(),
		Destroyed:   secret.VersionMetadata.Destroyed,
	}

	if !secret.VersionMetadata.DeletionTime.IsZero() {
		metadata.DeletionTime = secret.VersionMetadata.DeletionTime.String()
	}

	if secret.CustomMetadata != nil {
		metadata.CustomMetadata = convertToStringMap(secret.CustomMetadata)
	}

	return &SecretData{
		Data:     secret.Data,
		Metadata: metadata,
	}, nil
}

func (v *SecureVault) PatchSecret(ctx context.Context, path string, data map[string]interface{}) (*SecretMetadata, error) {
	secretData := map[string]interface{}{
		"data": data,
	}

	secret, err := v.client.KVv2(getMountPath(path)).Patch(ctx, getSecretPath(path), secretData)
	if err != nil {
		return nil, errors.ErrInternal("failed to patch secret").WithError(err)
	}

	metadata := &SecretMetadata{
		Version:     secret.VersionMetadata.Version,
		CreatedTime: secret.VersionMetadata.CreatedTime.String(),
		Destroyed:   secret.VersionMetadata.Destroyed,
	}

	if secret.CustomMetadata != nil {
		metadata.CustomMetadata = convertToStringMap(secret.CustomMetadata)
	}

	return metadata, nil
}

func (v *SecureVault) DeleteSecret(ctx context.Context, path string) error {
	err := v.client.KVv2(getMountPath(path)).Delete(ctx, getSecretPath(path))
	if err != nil {
		return errors.ErrInternal("failed to delete secret").WithError(err)
	}
	return nil
}

func (v *SecureVault) DeleteSecretVersions(ctx context.Context, path string, versions []int) error {
	err := v.client.KVv2(getMountPath(path)).DeleteVersions(ctx, getSecretPath(path), versions)
	if err != nil {
		return errors.ErrInternal("failed to delete secret versions").WithError(err)
	}
	return nil
}

func (v *SecureVault) UndeleteSecretVersions(ctx context.Context, path string, versions []int) error {
	err := v.client.KVv2(getMountPath(path)).Undelete(ctx, getSecretPath(path), versions)
	if err != nil {
		return errors.ErrInternal("failed to undelete secret versions").WithError(err)
	}
	return nil
}

func (v *SecureVault) DestroySecretVersions(ctx context.Context, path string, versions []int) error {
	err := v.client.KVv2(getMountPath(path)).Destroy(ctx, getSecretPath(path), versions)
	if err != nil {
		return errors.ErrInternal("failed to destroy secret versions").WithError(err)
	}
	return nil
}

func (v *SecureVault) DeleteSecretMetadata(ctx context.Context, path string) error {
	err := v.client.KVv2(getMountPath(path)).DeleteMetadata(ctx, getSecretPath(path))
	if err != nil {
		return errors.ErrInternal("failed to delete secret metadata").WithError(err)
	}
	return nil
}

func (v *SecureVault) ListSecrets(ctx context.Context, path string) ([]string, error) {
	secret, err := v.client.Logical().ListWithContext(ctx, fmt.Sprintf("%s/metadata/%s", getMountPath(path), getSecretPath(path)))
	if err != nil {
		return nil, errors.ErrInternal("failed to list secrets").WithError(err)
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

func (v *SecureVault) WriteSecretTyped(ctx context.Context, path string, value interface{}) (*SecretMetadata, error) {
	data, err := structToMap(value)
	if err != nil {
		return nil, errors.ErrInternal("failed to convert value to map").WithError(err)
	}
	return v.WriteSecret(ctx, path, data)
}

func (v *SecureVault) ReadSecretTyped(ctx context.Context, path string, target interface{}) error {
	data, err := v.ReadSecret(ctx, path)
	if err != nil {
		return err
	}

	return mapToStruct(data, target)
}

func getMountPath(path string) string {

	return "secret"
}

func getSecretPath(path string) string {
	return path
}

func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func mapToStruct(m map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func convertToStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}
