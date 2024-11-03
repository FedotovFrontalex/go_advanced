package request

import (
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, req *http.Request) (*T, error) {

	body, err := Decode[T](req.Body)

	if err != nil {
		return nil, err
	}

	err = Validate(body)

	if err != nil {
		return nil, err
	}

	return &body, nil
}
