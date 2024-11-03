package email

import (
	"encoding/json"
	"errors"
	"validationApi/pkg/fileStorage"
)

func GetVerificationData() (*VerifyData, error) {
	currentData, err := fileStorage.Read()

	if err != nil {
		return nil, err
	}

	var verifyData VerifyData
	err = json.Unmarshal(currentData, &verifyData)

	if err != nil {
		return nil, errors.New("Can't parse json")
	}

	return &verifyData, nil
}
