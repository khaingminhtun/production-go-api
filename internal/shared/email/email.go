package email

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sender interface {
	Send(
		ctx context.Context,
		to string,
		subject string,
		htmlBody string,
	) error
}

type SendGridSender struct {
	client    *sendgrid.Client
	fromEmail string
	fromName  string
}

func NewSendGridSender(
	apiKey string,
	fromEmail string,
	fromName string,
) Sender {
	return &SendGridSender{
		client:    sendgrid.NewSendClient(apiKey),
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *SendGridSender) Send(
	ctx context.Context,
	to string,
	subject string,
	htmlBody string,
) error {

	from := mail.NewEmail(
		s.fromName,
		s.fromEmail,
	)

	toEmail := mail.NewEmail(
		"",
		to,
	)

	message := mail.NewSingleEmail(
		from,
		subject,
		toEmail,
		"",
		htmlBody,
	)

	response, err := s.client.SendWithContext(ctx, message)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf(
			"sendgrid returned status %d: %s",
			response.StatusCode,
			response.Body,
		)
	}

	return nil
}
