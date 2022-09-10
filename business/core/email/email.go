// Package email provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/dudakovict/social-network/business/data/email"
	"go.uber.org/zap"
)

type EmailServer struct {
	log     *zap.SugaredLogger
	address string
	sender  string
	auth    smtp.Auth
	email.UnimplementedEmailServer
}

func NewEmailServer(log *zap.SugaredLogger, address string, sender string, auth smtp.Auth) EmailServer {
	return EmailServer{
		log:     log,
		address: address,
		sender:  sender,
		auth:    auth,
	}
}

func (es *EmailServer) Send(ctx context.Context, req *email.EmailRequest) (*email.EmailResponse, error) {

	message := []byte("Subject: Social network\n" + "Welcome to my social network!")

	err := smtp.SendMail(es.address, es.auth, es.sender, []string{req.Email}, message)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	res := &email.EmailResponse{
		Message: "Successfuly sent an email to " + req.Email,
	}

	return res, nil
}
