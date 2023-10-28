package model

import (
	"fmt"
)

type ErrNotFound struct {
	message string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s", e.message)
}
