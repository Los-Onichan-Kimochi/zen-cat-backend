package utils

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendPINBySMS(toPhone string, pin string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER") // tu número virtual Twilio (ej: +15005550006)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Mensaje a enviar
	body := fmt.Sprintf("Tu código de recuperación AstroCat es: %s", pin)

	params := &openapi.CreateMessageParams{}
	params.SetTo(toPhone)     // Destinatario: debe estar en formato +51XXXXXXXXX
	params.SetFrom(fromPhone) // Desde tu número Twilio
	params.SetBody(body)

	_, err := client.ApiV2010.CreateMessage(params)
	return err
}
