package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	googleoauth2 "golang.org/x/oauth2/google"
)

const (
	oauthStateTTL       = 10 * time.Minute
	googleUserInfoURL   = "https://www.googleapis.com/oauth2/v2/userinfo"
)

func (uc *usecase) googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     uc.Config.GoogleOAuth.ClientID,
		ClientSecret: uc.Config.GoogleOAuth.ClientSecret,
		RedirectURL:  uc.Config.GoogleOAuth.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     googleoauth2.Endpoint,
	}
}

func (uc *usecase) GetGoogleAuthURL(ctx context.Context) (*GoogleAuthURLResponse, error) {
	if !uc.Config.GoogleOAuth.IsEnabled() {
		return nil, errors.New("GOOGLE_OAUTH_DISABLED", "Google OAuth is not configured", http.StatusServiceUnavailable)
	}

	state, err := generateOAuthState()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OAuth state").WithError(err)
	}

	if err := uc.InMemoryStore.StoreOAuthState(ctx, state, oauthStateTTL); err != nil {
		return nil, errors.ErrInternal("failed to store OAuth state").WithError(err)
	}

	oauthCfg := uc.googleOAuthConfig()
	authURL := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))

	return &GoogleAuthURLResponse{AuthURL: authURL}, nil
}

func (uc *usecase) HandleGoogleCallback(ctx context.Context, req *GoogleCallbackRequest) (*GoogleCallbackResponse, error) {
	if !uc.Config.GoogleOAuth.IsEnabled() {
		return nil, errors.New("GOOGLE_OAUTH_DISABLED", "Google OAuth is not configured", http.StatusServiceUnavailable)
	}

	valid, err := uc.InMemoryStore.GetAndDeleteOAuthState(ctx, req.State)
	if err != nil {
		return nil, errors.ErrInternal("failed to validate OAuth state").WithError(err)
	}
	if !valid {
		return nil, errors.New("INVALID_STATE", "Invalid or expired OAuth state", http.StatusBadRequest)
	}

	oauthCfg := uc.googleOAuthConfig()
	token, err := oauthCfg.Exchange(ctx, req.Code)
	if err != nil {
		return nil, errors.New("OAUTH_EXCHANGE_FAILED", "Failed to exchange authorization code", http.StatusBadRequest)
	}

	googleUser, err := fetchGoogleUserInfo(ctx, oauthCfg, token)
	if err != nil {
		return nil, errors.ErrInternal("failed to fetch Google user info").WithError(err)
	}

	if !googleUser.EmailVerified {
		return nil, errors.New("EMAIL_NOT_VERIFIED", "Google account email is not verified", http.StatusForbidden)
	}

	existingAuth, err := uc.UserAuthMethodRepo.GetByCredentialField(ctx, string(entity.AuthMethodGoogle), "google_id", googleUser.GoogleID)
	if err == nil && existingAuth != nil {
		return uc.loginExistingGoogleUser(ctx, existingAuth, googleUser, req)
	}

	email := strings.ToLower(googleUser.Email)
	existingUser, err := uc.UserRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return uc.linkGoogleToExistingUser(ctx, existingUser, googleUser, req)
	}

	return uc.registerNewGoogleUser(ctx, googleUser, req)
}

func (uc *usecase) loginExistingGoogleUser(
	ctx context.Context,
	authMethod *entity.UserAuthMethod,
	googleUser *googleUserInfo,
	req *GoogleCallbackRequest,
) (*GoogleCallbackResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, authMethod.UserID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}

	if !user.IsActive() {
		return nil, errors.New("ACCOUNT_INACTIVE", "Your account is not active", http.StatusForbidden)
	}

	credData := entity.GoogleCredentialData{
		GoogleID:      googleUser.GoogleID,
		Email:         googleUser.Email,
		EmailVerified: googleUser.EmailVerified,
		Name:          googleUser.Name,
		Picture:       googleUser.Picture,
	}
	credJSON, err := json.Marshal(credData)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal credential data").WithError(err)
	}
	authMethod.CredentialData = credJSON
	authMethod.UpdatedAt = time.Now()
	if err := uc.UserAuthMethodRepo.Update(ctx, authMethod); err != nil {
		return nil, errors.ErrInternal("failed to update auth method").WithError(err)
	}

	return uc.issueGoogleTokens(ctx, user, req, false)
}

func (uc *usecase) linkGoogleToExistingUser(
	ctx context.Context,
	user *entity.User,
	googleUser *googleUserInfo,
	req *GoogleCallbackRequest,
) (*GoogleCallbackResponse, error) {
	if !user.IsActive() {
		return nil, errors.New("ACCOUNT_INACTIVE", "Your account is not active", http.StatusForbidden)
	}

	googleAuth := entity.NewGoogleAuthMethod(user.ID, entity.GoogleCredentialData{
		GoogleID:      googleUser.GoogleID,
		Email:         googleUser.Email,
		EmailVerified: googleUser.EmailVerified,
		Name:          googleUser.Name,
		Picture:       googleUser.Picture,
	})

	if err := uc.UserAuthMethodRepo.Create(ctx, googleAuth); err != nil {
		return nil, errors.ErrInternal("failed to create Google auth method").WithError(err)
	}

	secState, err := uc.UserSecurityStateRepo.GetByUserID(ctx, user.ID)
	if err == nil && secState != nil && !secState.EmailVerified {
		secState.EmailVerified = true
		now := time.Now()
		secState.EmailVerifiedAt = &now
		secState.UpdatedAt = now
		_ = uc.UserSecurityStateRepo.Update(ctx, secState)
	}

	return uc.issueGoogleTokens(ctx, user, req, false)
}

func (uc *usecase) registerNewGoogleUser(
	ctx context.Context,
	googleUser *googleUserInfo,
	req *GoogleCallbackRequest,
) (*GoogleCallbackResponse, error) {
	email := strings.ToLower(googleUser.Email)
	now := time.Now()

	user := &entity.User{
		Email:              email,
		Status:             entity.UserStatusActive,
		StatusChangedAt:    &now,
		RegistrationSource: "GOOGLE",
		Version:            1,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	var newUser *entity.User
	if err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.UserRepo.Create(txCtx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		newUser = user

		firstName, lastName := splitName(googleUser.Name)
		profile := &entity.UserProfile{
			UserID:            user.ID,
			FirstName:         firstName,
			LastName:          lastName,
			ProfilePictureURL: nilIfEmpty(googleUser.Picture),
			Metadata:          json.RawMessage("{}"),
			UpdatedAt:         now,
		}
		if err := uc.UserProfileRepo.Create(txCtx, profile); err != nil {
			return fmt.Errorf("create profile: %w", err)
		}

		googleAuth := entity.NewGoogleAuthMethod(user.ID, entity.GoogleCredentialData{
			GoogleID:      googleUser.GoogleID,
			Email:         googleUser.Email,
			EmailVerified: googleUser.EmailVerified,
			Name:          googleUser.Name,
			Picture:       googleUser.Picture,
		})
		if err := uc.UserAuthMethodRepo.Create(txCtx, googleAuth); err != nil {
			return fmt.Errorf("create auth method: %w", err)
		}

		secState := &entity.UserSecurityState{
			UserID:        user.ID,
			EmailVerified: true,
			EmailVerifiedAt: &now,
			UpdatedAt:     now,
		}
		if err := uc.UserSecurityStateRepo.Create(txCtx, secState); err != nil {
			return fmt.Errorf("create security state: %w", err)
		}

		return nil
	}); err != nil {
		return nil, errors.ErrInternal("failed to register Google user").WithError(err)
	}

	return uc.issueGoogleTokens(ctx, newUser, req, true)
}

func (uc *usecase) issueGoogleTokens(
	ctx context.Context,
	user *entity.User,
	req *GoogleCallbackRequest,
	isNewUser bool,
) (*GoogleCallbackResponse, error) {
	tenantClaims, userTenants, platformRoles, err := uc.buildMultiTenantClaims(ctx, user.ID)
	if err != nil {
		return nil, errors.ErrInternal("failed to build tenant claims").WithError(err)
	}

	sessionID := uuid.New()
	tokenFamily := uuid.New()

	tokenConfig, err := uc.buildTokenConfig()
	if err != nil {
		return nil, err
	}

	accessToken, err := jwtpkg.GenerateMultiTenantAccessToken(
		user.ID,
		user.Email,
		platformRoles,
		tenantClaims,
		sessionID,
		tokenConfig,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate access token").WithError(err)
	}

	refreshToken, err := jwtpkg.GenerateRefreshToken(user.ID, sessionID, tokenConfig)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	refreshTokenHash := hashToken(refreshToken)
	now := time.Now()
	refreshTokenEntity := &entity.RefreshToken{
		UserID:      user.ID,
		TokenHash:   refreshTokenHash,
		TokenFamily: tokenFamily,
		ExpiresAt:   now.Add(uc.Config.JWT.RefreshExpiry),
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		CreatedAt:   now,
	}

	userSession := &entity.UserSession{
		UserID:       user.ID,
		IPAddress:    req.IPAddress,
		LoginMethod:  entity.UserSessionLoginMethodGoogleOAuth,
		Status:       entity.UserSessionStatusActive,
		LastActiveAt: now,
		ExpiresAt:    now.Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if req.UserAgent != "" {
		userSession.UserAgent = &req.UserAgent
	}

	if err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.Create(txCtx, refreshTokenEntity); err != nil {
			return err
		}
		userSession.RefreshTokenID = &refreshTokenEntity.ID
		if err := uc.UserSessionRepo.Create(txCtx, userSession); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.ErrInternal("failed to create session").WithError(err)
	}

	now = time.Now()
	secState, err := uc.UserSecurityStateRepo.GetByUserID(ctx, user.ID)
	if err == nil && secState != nil {
		secState.LastLoginAt = &now
		secState.UpdatedAt = now
		_ = uc.UserSecurityStateRepo.Update(ctx, secState)
	}

	profile, _ := uc.UserProfileRepo.GetByUserID(ctx, user.ID)
	fullName := ""
	if profile != nil {
		fullName = profile.FirstName
		if profile.LastName != "" {
			fullName += " " + profile.LastName
		}
	}

	return &GoogleCallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(uc.Config.JWT.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
		IsNewUser:    isNewUser,
		User: LoginUserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: fullName,
			Tenants:  userTenants,
		},
	}, nil
}

type googleUserInfo struct {
	GoogleID      string `json:"id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*googleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	var userInfo googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	if userInfo.GoogleID == "" {
		return nil, fmt.Errorf("google user info missing id")
	}

	return &userInfo, nil
}

func generateOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func splitName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "User", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
