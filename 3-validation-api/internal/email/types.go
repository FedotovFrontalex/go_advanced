package email

type VerifyData struct {
		Data []EmailVerifyData `json:"data"`
}

type EmailVerifyData struct {
		Email string `json:"email"`
		Hash string `json:"hash"`
}
