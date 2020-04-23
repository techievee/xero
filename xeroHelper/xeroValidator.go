package xeroHelper

import (
	"github.com/google/uuid"
)

func ValidateUUID(s string) bool {
	_, err := uuid.Parse(s)
	if err != nil {
		return false
	} else {
		return true
	}
}
