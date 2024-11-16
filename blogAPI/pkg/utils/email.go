package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/anoying-kid/go-apps/blogAPI/pkg/config"
)

type EmailTemplate struct {
    Subject string
    Body    string
}

func SendPasswordResetEmail(email, resetToken string, config config.Config) error {
    // Create HTML template for the email
    templateStr := `
<!DOCTYPE html>
<html>
<head>
    <style>
        .container {
            padding: 20px;
            font-family: Arial, sans-serif;
        }
        .button {
            background-color: #4CAF50;
            border: none;
            color: white;
            padding: 15px 32px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 16px;
            margin: 4px 2px;
            cursor: pointer;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Password Reset Request</h2>
        <p>Hello,</p>
        <p>You have requested to reset your password. Click the button below to set a new password:</p>
        <p>
            <a href="{{.ResetLink}}" class="button">Reset Password</a>
        </p>
        <p>Or copy and paste this link in your browser:</p>
        <p>{{.ResetLink}}</p>
        <p>This link will expire in 1 hour.</p>
        <p>If you didn't request this, please ignore this email.</p>
        <br>
        <p>Best regards,<br>Your Application Team</p>
    </div>
</body>
</html>`

    // Parse the template
    t, err := template.New("resetEmail").Parse(templateStr)
    if err != nil {
        return fmt.Errorf("error parsing email template: %w", err)
    }

    // Prepare template data
    data := struct {
        ResetLink string
    }{
        ResetLink: fmt.Sprintf("%s/reset-password?token=%s", config.Frontend.URL, resetToken),
    }

    // Execute the template
    var body bytes.Buffer
    if err := t.Execute(&body, data); err != nil {
        return fmt.Errorf("error executing email template: %w", err)
    }

    // Prepare email headers and body
    headers := make(map[string]string)
    headers["From"] = config.Email.Username
    headers["To"] = email
    headers["Subject"] = "Password Reset Request"
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/html; charset=utf-8"

    // Construct message
    message := ""
    for k, v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body.String()

    // Authentication
    auth := smtp.PlainAuth(
        "",
        config.Email.Username,
        config.Email.Password,
        config.Email.Host,
    )

    // Send email
    err = smtp.SendMail(
        fmt.Sprintf("%s:%d", config.Email.Host, config.Email.Port),
        auth,
        config.Email.Username,
        []string{email},
        []byte(message),
    )

    if err != nil {
        return fmt.Errorf("error sending email: %w", err)
    }

    return nil
}