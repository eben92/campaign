package utils

import (
	"encoding/json"
	"log/slog"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WrapInResponse(message string, data interface{}) []byte {
	res := ApiResponse{
		Message: message,
		Data:    data,
	}
	b, err := json.Marshal(res)

	if err != nil {
		slog.Error("Error marshalling response", "error", err)

		return nil
	}

	return b
}
