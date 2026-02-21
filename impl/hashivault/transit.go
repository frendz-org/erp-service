package hashivault

import (
	"context"
	"encoding/base64"
	"erp-service/pkg/errors"
	"fmt"
)

type TransitKey struct {
	Name                 string           `json:"name"`
	Type                 string           `json:"type"`
	DeletionAllowed      bool             `json:"deletion_allowed"`
	Derived              bool             `json:"derived"`
	Exportable           bool             `json:"exportable"`
	AllowPlaintextBackup bool             `json:"allow_plaintext_backup"`
	Keys                 map[string]int64 `json:"keys"`
	MinDecryptionVersion int              `json:"min_decryption_version"`
	MinEncryptionVersion int              `json:"min_encryption_version"`
	SupportsEncryption   bool             `json:"supports_encryption"`
	SupportsDecryption   bool             `json:"supports_decryption"`
	SupportsSigning      bool             `json:"supports_signing"`
	SupportsDerivation   bool             `json:"supports_derivation"`
}

func (v *SecureVault) EncryptData(ctx context.Context, keyName string, plaintext []byte) (string, error) {
	encodedPlaintext := base64.StdEncoding.EncodeToString(plaintext)

	data := map[string]interface{}{
		"plaintext": encodedPlaintext,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/encrypt/%s", keyName), data)
	if err != nil {
		return "", errors.ErrInternal("failed to encrypt data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected ciphertext format")
	}

	return ciphertext, nil
}

func (v *SecureVault) DecryptData(ctx context.Context, keyName string, ciphertext string) ([]byte, error) {
	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/decrypt/%s", keyName), data)
	if err != nil {
		return nil, errors.ErrInternal("failed to decrypt data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("vault returned empty response")
	}

	encodedPlaintext, ok := secret.Data["plaintext"].(string)
	if !ok {
		return nil, errors.ErrInternal("unexpected plaintext format")
	}

	plaintext, err := base64.StdEncoding.DecodeString(encodedPlaintext)
	if err != nil {
		return nil, errors.ErrInternal("failed to decode plaintext").WithError(err)
	}

	return plaintext, nil
}

func (v *SecureVault) EncryptDataWithContext(ctx context.Context, keyName string, plaintext []byte, keyContext []byte) (string, error) {
	encodedPlaintext := base64.StdEncoding.EncodeToString(plaintext)
	encodedContext := base64.StdEncoding.EncodeToString(keyContext)

	data := map[string]interface{}{
		"plaintext": encodedPlaintext,
		"context":   encodedContext,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/encrypt/%s", keyName), data)
	if err != nil {
		return "", errors.ErrInternal("failed to encrypt data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected ciphertext format")
	}

	return ciphertext, nil
}

func (v *SecureVault) DecryptDataWithContext(ctx context.Context, keyName string, ciphertext string, keyContext []byte) ([]byte, error) {
	encodedContext := base64.StdEncoding.EncodeToString(keyContext)

	data := map[string]interface{}{
		"ciphertext": ciphertext,
		"context":    encodedContext,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/decrypt/%s", keyName), data)
	if err != nil {
		return nil, errors.ErrInternal("failed to decrypt data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("vault returned empty response")
	}

	encodedPlaintext, ok := secret.Data["plaintext"].(string)
	if !ok {
		return nil, errors.ErrInternal("unexpected plaintext format")
	}

	plaintext, err := base64.StdEncoding.DecodeString(encodedPlaintext)
	if err != nil {
		return nil, errors.ErrInternal("failed to decode plaintext").WithError(err)
	}

	return plaintext, nil
}

func (v *SecureVault) RewrapData(ctx context.Context, keyName string, ciphertext string) (string, error) {
	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/rewrap/%s", keyName), data)
	if err != nil {
		return "", errors.ErrInternal("failed to rewrap data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	newCiphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected ciphertext format")
	}

	return newCiphertext, nil
}

func (v *SecureVault) GenerateDataKey(ctx context.Context, keyName string, bits int) (string, string, error) {
	data := map[string]interface{}{}
	if bits > 0 {
		data["bits"] = bits
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/datakey/plaintext/%s", keyName), data)
	if err != nil {
		return "", "", errors.ErrInternal("failed to generate data key").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", "", errors.ErrNotFound("vault returned empty response")
	}

	plaintext, ok := secret.Data["plaintext"].(string)
	if !ok {
		return "", "", errors.ErrInternal("unexpected plaintext format")
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", "", errors.ErrInternal("unexpected ciphertext format")
	}

	return plaintext, ciphertext, nil
}

func (v *SecureVault) GenerateWrappedDataKey(ctx context.Context, keyName string, bits int) (string, error) {
	data := map[string]interface{}{}
	if bits > 0 {
		data["bits"] = bits
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/datakey/wrapped/%s", keyName), data)
	if err != nil {
		return "", errors.ErrInternal("failed to generate wrapped data key").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected ciphertext format")
	}

	return ciphertext, nil
}

func (v *SecureVault) CreateTransitKey(ctx context.Context, keyName string, keyType string, derivable bool) error {
	data := map[string]interface{}{
		"type":    keyType,
		"derived": derivable,
	}

	_, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/keys/%s", keyName), data)
	if err != nil {
		return errors.ErrInternal("failed to create transit key").WithError(err)
	}

	return nil
}

func (v *SecureVault) ReadTransitKey(ctx context.Context, keyName string) (*TransitKey, error) {
	secret, err := v.client.Logical().ReadWithContext(ctx, fmt.Sprintf("transit/keys/%s", keyName))
	if err != nil {
		return nil, errors.ErrInternal("failed to read transit key").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.ErrNotFound("transit key not found")
	}

	transitKey := &TransitKey{
		Name: keyName,
	}

	if v, ok := secret.Data["type"].(string); ok {
		transitKey.Type = v
	}
	if v, ok := secret.Data["deletion_allowed"].(bool); ok {
		transitKey.DeletionAllowed = v
	}
	if v, ok := secret.Data["derived"].(bool); ok {
		transitKey.Derived = v
	}
	if v, ok := secret.Data["exportable"].(bool); ok {
		transitKey.Exportable = v
	}
	if v, ok := secret.Data["allow_plaintext_backup"].(bool); ok {
		transitKey.AllowPlaintextBackup = v
	}
	if v, ok := secret.Data["supports_encryption"].(bool); ok {
		transitKey.SupportsEncryption = v
	}
	if v, ok := secret.Data["supports_decryption"].(bool); ok {
		transitKey.SupportsDecryption = v
	}
	if v, ok := secret.Data["supports_signing"].(bool); ok {
		transitKey.SupportsSigning = v
	}
	if v, ok := secret.Data["supports_derivation"].(bool); ok {
		transitKey.SupportsDerivation = v
	}

	return transitKey, nil
}

func (v *SecureVault) DeleteTransitKey(ctx context.Context, keyName string) error {
	_, err := v.client.Logical().DeleteWithContext(ctx, fmt.Sprintf("transit/keys/%s", keyName))
	if err != nil {
		return errors.ErrInternal("failed to delete transit key").WithError(err)
	}
	return nil
}

func (v *SecureVault) RotateTransitKey(ctx context.Context, keyName string) error {
	_, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/keys/%s/rotate", keyName), nil)
	if err != nil {
		return errors.ErrInternal("failed to rotate transit key").WithError(err)
	}
	return nil
}

func (v *SecureVault) SignData(ctx context.Context, keyName string, input []byte) (string, error) {
	encodedInput := base64.StdEncoding.EncodeToString(input)

	data := map[string]interface{}{
		"input": encodedInput,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/sign/%s", keyName), data)
	if err != nil {
		return "", errors.ErrInternal("failed to sign data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	signature, ok := secret.Data["signature"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected signature format")
	}

	return signature, nil
}

func (v *SecureVault) VerifySignature(ctx context.Context, keyName string, input []byte, signature string) (bool, error) {
	encodedInput := base64.StdEncoding.EncodeToString(input)

	data := map[string]interface{}{
		"input":     encodedInput,
		"signature": signature,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/verify/%s", keyName), data)
	if err != nil {
		return false, errors.ErrInternal("failed to verify signature").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return false, errors.ErrNotFound("vault returned empty response")
	}

	valid, ok := secret.Data["valid"].(bool)
	if !ok {
		return false, errors.ErrInternal("unexpected valid format")
	}

	return valid, nil
}

func (v *SecureVault) GenerateRandomBytes(ctx context.Context, numBytes int, format string) (string, error) {
	data := map[string]interface{}{
		"bytes": numBytes,
	}

	if format != "" {
		data["format"] = format
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "transit/random", data)
	if err != nil {
		return "", errors.ErrInternal("failed to generate random bytes").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	randomBytes, ok := secret.Data["random_bytes"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected random_bytes format")
	}

	return randomBytes, nil
}

func (v *SecureVault) Hash(ctx context.Context, input []byte, algorithm string) (string, error) {
	encodedInput := base64.StdEncoding.EncodeToString(input)

	data := map[string]interface{}{
		"input": encodedInput,
	}

	if algorithm != "" {
		data["algorithm"] = algorithm
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, "transit/hash", data)
	if err != nil {
		return "", errors.ErrInternal("failed to hash data").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	sum, ok := secret.Data["sum"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected sum format")
	}

	return sum, nil
}

func (v *SecureVault) GenerateHMAC(ctx context.Context, keyName string, input []byte) (string, error) {
	encodedInput := base64.StdEncoding.EncodeToString(input)

	data := map[string]interface{}{
		"input": encodedInput,
	}

	secret, err := v.client.Logical().WriteWithContext(ctx, fmt.Sprintf("transit/hmac/%s", keyName), data)
	if err != nil {
		return "", errors.ErrInternal("failed to generate HMAC").WithError(err)
	}

	if secret == nil || secret.Data == nil {
		return "", errors.ErrNotFound("vault returned empty response")
	}

	hmac, ok := secret.Data["hmac"].(string)
	if !ok {
		return "", errors.ErrInternal("unexpected hmac format")
	}

	return hmac, nil
}
