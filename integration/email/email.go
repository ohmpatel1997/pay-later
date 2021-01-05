package email

import (
	"pay-later/integration/log"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type EmailService interface {
	IsValid(string) bool
}

type emailService struct {
	log log.Logger
}

func NewEmailService(l log.Logger) EmailService {
	return &emailService{
		l,
	}
}

func (e emailService) IsValid(email string) bool {
	if len(email) < 3 && len(email) > 254 { //max length of any email
		return false
	}
	if !emailRegex.MatchString(email) {
		return false
	}

	return true
}
