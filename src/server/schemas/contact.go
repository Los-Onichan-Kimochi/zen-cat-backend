package schemas

type ContactRequest struct {
	Name    string `json:"name" `
	Email   string `json:"email" `
	Phone   string `json:"phone,omitempty" ` // opcional
	Subject string `json:"subject" `         // requerido
	Message string `json:"message" `
}
