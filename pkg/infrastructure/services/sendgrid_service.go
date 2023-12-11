package services

import (
	"encoding/json"
	"fmt"
	"github.com/ndodanli/go-clean-architecture/configs"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

type ISendgridService interface {
	SendEmail(to string, subject string, htmlContent string) error
}

type SendgridService struct {
	sendgridCfg *configs.Sendgrid
	client      *sendgrid.Client
}

func NewSendgridService(sendgridCfg *configs.Sendgrid) ISendgridService {
	return &SendgridService{
		sendgridCfg: sendgridCfg,
		client:      sendgrid.NewSendClient(sendgridCfg.API_KEY),
	}
}

func (s *SendgridService) SendEmail(to string, subject string, htmlContent string) error {
	fromEmail := mail.NewEmail(s.sendgridCfg.FROM_NAME, s.sendgridCfg.FROM_EMAIL)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(fromEmail, subject, toEmail, "", htmlContent)
	response, err := s.client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	if response.StatusCode != 202 {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
		var jsonResponse []byte
		jsonResponse, err = json.Marshal(response)
		if err != nil {
			return err
		}
		return httperr.SendgridError(string(jsonResponse))
	}

	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)

	return nil
}
