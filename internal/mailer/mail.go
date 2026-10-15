package mailer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailer() (*Mailer, error) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return &Mailer{}, err
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))

	return &Mailer{
		dialer: d,
		from:   d.Username,
	}, nil
}

func RenderTemplate(filename string, data any) (string, error) {
	path := filepath.Join("internal/template", filename)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *Mailer) SendMail(to, html string) error {
	mes := gomail.NewMessage()
	mes.SetAddressHeader("From", m.from, "E-Wallet")
	mes.SetHeader("To", to)
	mes.SetHeader("Subject", "OTP for Reset Password")
	mes.SetBody("text/html", html)

	return m.dialer.DialAndSend(mes)
}
