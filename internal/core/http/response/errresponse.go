package httpresp

import (
	"errors"
	"net/http"
	"time"

	"github.com/dirijable/coworking-api/internal/core/error/apierror"
	"github.com/dirijable/coworking-api/internal/core/error/apperror"
)

func SendErrorResponse(rw http.ResponseWriter, err error) {
	var (
		statusCode int
		msg        string
		details    = make(map[string]string)
	)
	switch {
	case errors.Is(err, apperror.ErrBadRequest):
		msg = "Validation error"
		statusCode = http.StatusBadRequest
		if ve, ok := errors.AsType[apperror.ValidationError](err); ok {
			details = ve.Fields
		}
	case errors.Is(err, apperror.ErrConflict):
		msg = "Resource already exist"
		statusCode = http.StatusConflict
	case errors.Is(err, apperror.ErrNotFound):
		msg = "Resource not found"
		statusCode = http.StatusNotFound
	default:
		msg = "Internal server error"
		statusCode = http.StatusInternalServerError
	}
	sendError(rw, msg, statusCode, details)
}

func sendError(rw http.ResponseWriter, msg string, statusCode int, details map[string]string) {
	apiError := apierror.APIError{
		Msg:       msg,
		Timestamp: time.Now(),
		Details:   details,
	}
	_ = SendJSONResponse(rw, statusCode, apiError)
}
