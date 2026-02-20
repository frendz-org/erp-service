package entity

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordAuthMethod(t *testing.T) {
	userID := uuid.New()
	passwordHash := "$2a$10$examplehashedpassword"

	t.Run("returns valid UserAuthMethod with correct fields", func(t *testing.T) {
		method := NewPasswordAuthMethod(userID, passwordHash)

		require.NotNil(t, method)
		assert.Equal(t, userID, method.UserID)
		assert.Equal(t, string(AuthMethodPassword), method.MethodType)
		assert.True(t, method.IsActive)
	})

	t.Run("CredentialData is non-empty valid JSON", func(t *testing.T) {
		method := NewPasswordAuthMethod(userID, passwordHash)

		require.NotEmpty(t, method.CredentialData,
			"CredentialData must not be empty — silent json.Marshal error would produce nil/empty bytes")

		var data PasswordCredentialData
		err := json.Unmarshal(method.CredentialData, &data)
		require.NoError(t, err, "CredentialData must be valid JSON")
	})

	t.Run("CredentialData contains the provided password hash", func(t *testing.T) {
		method := NewPasswordAuthMethod(userID, passwordHash)

		var data PasswordCredentialData
		require.NoError(t, json.Unmarshal(method.CredentialData, &data))
		assert.Equal(t, passwordHash, data.PasswordHash)
	})

	t.Run("CredentialData contains zero password history entries", func(t *testing.T) {
		method := NewPasswordAuthMethod(userID, passwordHash)

		var data PasswordCredentialData
		require.NoError(t, json.Unmarshal(method.CredentialData, &data))

		assert.Empty(t, data.PasswordHistory, "PasswordHistory should contain no entries")
	})

	t.Run("GetPasswordHash returns the hash stored in CredentialData", func(t *testing.T) {
		method := NewPasswordAuthMethod(userID, passwordHash)
		assert.Equal(t, passwordHash, method.GetPasswordHash())
	})
}
