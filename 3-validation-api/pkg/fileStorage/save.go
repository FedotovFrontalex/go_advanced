package fileStorage

import (
	"errors"
	"os"
)

func Save(content []byte) error {
		file, err := os.Create("storage.json")

		if err != nil {
				return errors.New("Unable to create file: " + err.Error())
		}

		defer file.Close()

		_, err = file.Write(content)

		if err != nil {
				return errors.New("Unable to write file: " + err.Error())
		}

		return nil
}
