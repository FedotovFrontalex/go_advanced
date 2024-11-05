package fileStorage

import (
	"errors"
	"os"
)

func Read() ([]byte, error) {
	data, err := os.ReadFile("storage.json")

	if err != nil {	
		return nil, errors.New("Unable to read file: " + err.Error())
	}

	return data, nil
}
