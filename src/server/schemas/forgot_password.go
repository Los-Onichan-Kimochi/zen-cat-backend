package schemas

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Pin     string `json:"pin"`
}
