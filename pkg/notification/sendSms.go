package notification

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	"go-ecommerce-app/config"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type NotificationClient interface {
	SendSMS(phone string, message string) error
}

type notificationClient struct {
	config config.AppConfig
}

// Twilio
func (c notificationClient) SendSMS(phone string, message string) error {

	accountSid := c.config.TwilioAccountSid
	authToken := c.config.TwilioAuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo("+917615942926")                  // Mynumber
	params.SetFrom(c.config.TwilioFromPhoneNumber) // from twillio
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

	return nil
}

func NewNotificationClient(cfg config.AppConfig) NotificationClient {
	return &notificationClient{
		config: cfg,
	}
}
