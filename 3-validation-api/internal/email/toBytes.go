package email

import (
	"encoding/json"
	"errors"
)

func ToBytes(data *VerifyData) ([]byte, error) {
		bytes, err := json.Marshal(data)

		if err != nil {
				return nil, errors.New("Can't convert data to json")
		}

		return bytes, nil
}
