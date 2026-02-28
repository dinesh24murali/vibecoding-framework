package notifications

import (
	"context"
	"fmt"
	"log"
	"regexp"
)

var phonePattern = regexp.MustCompile(`^[6-9][0-9]{9}$`)

type SMSClient interface {
	Send(ctx context.Context, phoneNumber string, message string) error
}

type LogSMSClient struct{}

func NewLogSMSClient() *LogSMSClient {
	return &LogSMSClient{}
}

func (c *LogSMSClient) Send(_ context.Context, phoneNumber string, message string) error {
	if !phonePattern.MatchString(phoneNumber) {
		return fmt.Errorf("invalid phone number")
	}

	log.Printf("sms dispatched phone=%s message=%q", phoneNumber, message)
	return nil
}
