package response

import (
	"encoding/json"
	"net/http"
	"validationApi/pkg/logger"
)

func Json(w http.ResponseWriter, res any, status int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
				logger.Error(err)
		}
}
