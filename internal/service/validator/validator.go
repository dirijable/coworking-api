package validator

import (
	"errors"
	"strings"

	"github.com/dirijable/coworking-api/internal/domain"
)

func Validate(room domain.Room) error {
	if strings.TrimSpace(room.Name) == "" {
		return errors.New("empty name")
	}
	if room.Capacity <= 0 {
		return errors.New("capacity <= 0")
	}
	return nil
}
