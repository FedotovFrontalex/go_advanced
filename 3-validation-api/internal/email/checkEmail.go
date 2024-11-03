package email

import (
	"errors"
)

func CheckEmail(hash string) (int, error) {
	data, err := GetVerificationData()

	if err != nil {
		return 500, err
	}

	var noVerified []EmailVerifyData

	for _, val := range data.Data {
		if val.Hash != hash {
			noVerified = append(noVerified, val)
		}
	}

	if len(data.Data) == len(noVerified) {
		return 404, errors.New("Not found")
	}

	data.Data = noVerified
		
	err = SaveData(data)

	if err != nil {
		return 500, err
	}
		
	return 200, nil
}
