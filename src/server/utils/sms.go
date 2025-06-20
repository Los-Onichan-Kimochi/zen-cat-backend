package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendPINBySMS envía un PIN a través de SMS utilizando Twilio
func SendPINBySMS(toPhone string, pin string) error {
	// Obtener las credenciales desde las variables de entorno
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER") // Tu número virtual de Twilio (ej: +15005550006)

	// Crear un cliente de Twilio usando las credenciales
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Crear el cuerpo del mensaje
	body := fmt.Sprintf("Tu código de recuperación AstroCat es: %s", pin)

	// Crear los parámetros para el mensaje
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(toPhone)     // Destinatario: debe estar en formato +51XXXXXXXXX
	params.SetFrom(fromPhone) // Desde tu número Twilio
	params.SetBody(body)      // Contenido del mensaje

	// Enviar el mensaje
	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
		return err
	}

	// Imprimir la respuesta de la API (puedes omitirlo en producción)
	response, _ := json.Marshal(*resp)
	fmt.Println("Response: " + string(response))
	return nil
}
