package email

func SaveVerifyEmailData(data EmailVerifyData) error {
	verifyData, _ := GetVerificationData()

	if verifyData == nil {
		verifyData = &VerifyData{
			Data: []EmailVerifyData{},
		}
	}

	var emailsData []EmailVerifyData

	for _, val := range verifyData.Data {
		if val.Email != data.Email {
			emailsData = append(emailsData, val)
		}
	}

	emailsData = append(emailsData, data)

	verifyData.Data = emailsData

	return SaveData(verifyData)
}
