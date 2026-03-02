package auth

const (
	PasswordMinLength = 8
)

const (
	OTPLength = 6
)

const (
	RegistrationSessionExpiryMinutes = 10

	RegistrationOTPLength         = 6
	RegistrationOTPExpiryMinutes  = 10
	RegistrationOTPMaxAttempts    = 5
	RegistrationOTPMaxResends     = 3
	RegistrationOTPResendCooldown = 60

	RegistrationCompleteTokenExpiryMinutes = 5
	RegistrationCompleteTokenPurpose       = "registration_complete"

	RegistrationRateLimitPerHour = 3
	RegistrationRateLimitWindow  = 60
)

const (
	LoginSessionExpiryMinutes = 10
	LoginOTPExpiryMinutes     = 5
	LoginOTPMaxAttempts       = 5
	LoginOTPMaxResends        = 3
	LoginOTPResendCooldown    = 60
	LoginRateLimitPerHour     = 5
	LoginRateLimitWindow      = 60
)

const (
	BcryptTargetCost = 12
)

const (
	TransferTokenTTLSeconds      = 30
	TransferTokenCodeBytes       = 32
	TransferTokenRateLimitPerMin = 5
	TransferTokenRateLimitWindow = 1
	TransferTokenMaxTreeDepth    = 20
)
