package common

import (
	"fmt"
	"io"
)

func DecodeError(body io.ReadCloser) error {
	errorBody, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("Cannot even decode body for error")
	}
	return fmt.Errorf("Status code different than 200, %s", string(errorBody))
}
