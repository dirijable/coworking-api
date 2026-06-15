package httpresp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJSONResponse(rw http.ResponseWriter, statusCode int, body any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	if err := json.NewEncoder(rw).Encode(body); err != nil {
		return fmt.Errorf("body encode: %w", err)
	}
	return nil
}
