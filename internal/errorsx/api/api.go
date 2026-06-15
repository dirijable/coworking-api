package api

import (
	"time"
)

type APIError struct {
	Msg       string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
}
