package email

import "validationApi/pkg/fileStorage"

func SaveData(data *VerifyData) error {
	bytes, err := ToBytes(data)

	if err != nil {
		return err
	}

	err = fileStorage.Save(bytes)

	if err != nil {
		return err
	}

	return nil
}
