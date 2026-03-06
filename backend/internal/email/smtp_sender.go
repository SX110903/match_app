package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
	"runtime"

	"github.com/SX110903/match_app/backend/internal/config"
)

type smtpSender struct {
	cfg config.EmailConfig
}

func NewSMTPSender(cfg config.EmailConfig) IEmailService {
	return &smtpSender{cfg: cfg}
}

func (s *smtpSender) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.SMTPHost)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		s.cfg.From, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.User, []string{to}, []byte(msg))
}

func (s *smtpSender) renderTemplate(name string, data any) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	templatePath := filepath.Join(filepath.Dir(filename), "templates", name)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", name, err)
	}
	return buf.String(), nil
}

func (s *smtpSender) SendVerificationEmail(_ context.Context, to, name, verifyURL string) error {
	body, err := s.renderTemplate("verify_email.html", map[string]string{
		"Name": name, "URL": verifyURL,
	})
	if err != nil {
		return err
	}
	return s.send(to, "Verifica tu email - MatchHub", body)
}

func (s *smtpSender) SendPasswordResetEmail(_ context.Context, to, resetURL string) error {
	body, err := s.renderTemplate("reset_password.html", map[string]string{"URL": resetURL})
	if err != nil {
		return err
	}
	return s.send(to, "Restablecer contraseña - MatchHub", body)
}

func (s *smtpSender) SendPasswordChangedEmail(_ context.Context, to string) error {
	return s.send(to, "Tu contraseña fue cambiada - MatchHub",
		"<p>Tu contraseña fue cambiada exitosamente. Si no fuiste tú, contacta soporte inmediatamente.</p>")
}

func (s *smtpSender) SendWelcomeEmail(_ context.Context, to, name string) error {
	body, err := s.renderTemplate("welcome.html", map[string]string{"Name": name})
	if err != nil {
		return err
	}
	return s.send(to, "Bienvenido a MatchHub", body)
}
