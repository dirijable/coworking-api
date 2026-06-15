package extractor

import (
	"fmt"
	"net/http"

	"github.com/dirijable/coworking-api/internal/errorsx/service"
	"github.com/google/uuid"
)

func PathUUID(r *http.Request, key string) (uuid.UUID, error) {
	valueStr := r.PathValue(key)
	value, err := uuid.Parse(valueStr)
	if err != nil {
		err = fmt.Errorf("%w: invalid path parameter: %v", service.ErrBadRequest, err)
		return uuid.Nil, err
	}
	return value, nil
}
