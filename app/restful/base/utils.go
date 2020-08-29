package base

import (
	"context"
	"encoding/json"
	"go-cicd/app/domain/entity"
	"go-cicd/app/domain/model"
	"net/http"
)

// ResponseError response error to client
func ResponseError(code int, message string, w http.ResponseWriter) {
	response := model.APIResponse{
		ErrorCode: code,
		Message:   message,
	}

	ResponseToClient(&response, w)
}

// ResponseToClient send response to client
func ResponseToClient(response *model.APIResponse, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	// w.Header().Add("Cache-Control", "max-age=60")
	json.NewEncoder(w).Encode(response)
}

// GetUserIDFromContext get user ID from context after authenticate
func GetUserIDFromContext(c context.Context) string {
	return GetUserFromContext(c).ID
}

// GetUserFromContext get user from context after authenticate
func GetUserFromContext(c context.Context) *entity.User {
	if user, ok := c.Value(UserAuthenticatedKey).(*entity.User); ok {
		return user
	}

	return nil
}
