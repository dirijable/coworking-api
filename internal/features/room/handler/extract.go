package handler

import (
	"fmt"
	"net/http"

	"github.com/dirijable/coworking-api/internal/core/error/apperror"
	"github.com/google/uuid"
)

func pathUUID(r *http.Request, key string) (uuid.UUID, error) {
	valueStr := r.PathValue(key)
	value, err := uuid.Parse(valueStr)
	if err != nil {
		err = fmt.Errorf("%w: invalid path parameter: %v", apperror.ErrBadRequest, err)
		return uuid.Nil, err
	}
	return value, nil
}
