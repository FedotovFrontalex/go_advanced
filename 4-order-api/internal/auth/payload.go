package auth

type AuthRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type AuthResponse struct {
	SessionId string `json:"sessionId"`
}

type AuthVerifyRequest struct {
	SessionId string `json:"sessionId" validate:"required"`
	Code      int    `json:"code" validate:"required"`
}

type AuthVerificationResponse struct {
	Token string `json:"token"`
}
