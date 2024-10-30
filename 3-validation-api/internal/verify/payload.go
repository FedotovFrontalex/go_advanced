package verify

type SendPayload struct {
		Email string `json:"email"`
}

type VerifyPayload struct {
		Email string `json:"email"`
		Address string `json:"address"`
}
