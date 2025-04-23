package commons

import (
	"encoding/json"
)

type ApiErrorResponsePayload struct {
	Status         string                 `json:"status"`
	Message        string                 `json:"message"`
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
}

func PrintStruct(payload interface{}) string {
	pbytes, _ := json.Marshal(payload)
	return string(pbytes)
}

func ApiErrorResponse(message string, additionalInfo map[string]interface{}) *ApiErrorResponsePayload {
	response := &ApiErrorResponsePayload{
		Status:  "Error",
		Message: message,
	}
	if len(additionalInfo) > 0 {
		response.AdditionalInfo = additionalInfo
	}
	return response
}
