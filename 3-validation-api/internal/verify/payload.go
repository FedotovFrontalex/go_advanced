package verify

type SendRequest struct {
		Email string `json:"email" validate:"required,email"`
}

type SendPayload struct {
		Email string `json:"email"`
}

type VerifyPayload struct {
		Email string `json:"email"`
		Address string `json:"address"`
}
