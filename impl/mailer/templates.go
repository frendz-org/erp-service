package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		panic(fmt.Sprintf("failed to parse email templates: %v", err))
	}
}

type OTPTemplateData struct {
	OTP           string
	ExpiryMinutes int
	Year          int
}

type WelcomeTemplateData struct {
	FirstName string
	Year      int
}

type PasswordResetTemplateData struct {
	Token         string
	ExpiryMinutes int
	Year          int
}

type PINResetTemplateData struct {
	OTP           string
	ExpiryMinutes int
	Year          int
}

type AdminInvitationTemplateData struct {
	Token         string
	ExpiryMinutes int
	Year          int
}

func renderTemplate(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}
	return buf.String(), nil
}

func renderRegistrationOTPEmail(otp string, expiryMinutes int) (string, error) {
	return renderTemplate("otp.html", OTPTemplateData{
		OTP:           otp,
		ExpiryMinutes: expiryMinutes,
		Year:          time.Now().Year(),
	})
}

func renderLoginOTPEmail(otp string, expiryMinutes int) (string, error) {
	return renderTemplate("login_otp.html", OTPTemplateData{
		OTP:           otp,
		ExpiryMinutes: expiryMinutes,
		Year:          time.Now().Year(),
	})
}

func renderWelcomeEmail(firstName string) (string, error) {
	return renderTemplate("welcome.html", WelcomeTemplateData{
		FirstName: firstName,
		Year:      time.Now().Year(),
	})
}

func renderPasswordResetEmail(token string, expiryMinutes int) (string, error) {
	return renderTemplate("password_reset.html", PasswordResetTemplateData{
		Token:         token,
		ExpiryMinutes: expiryMinutes,
		Year:          time.Now().Year(),
	})
}

func renderPINResetEmail(otp string, expiryMinutes int) (string, error) {
	return renderTemplate("pin_reset.html", PINResetTemplateData{
		OTP:           otp,
		ExpiryMinutes: expiryMinutes,
		Year:          time.Now().Year(),
	})
}

func renderAdminInvitationEmail(token string, expiryMinutes int) (string, error) {
	return renderTemplate("admin_invitation.html", AdminInvitationTemplateData{
		Token:         token,
		ExpiryMinutes: expiryMinutes,
		Year:          time.Now().Year(),
	})
}
