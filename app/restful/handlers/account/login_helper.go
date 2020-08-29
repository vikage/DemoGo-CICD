package account

import (
	"go-cicd/app/authenticate"
	"go-cicd/app/domain/entity"
	"go-cicd/app/domain/model"
	"go-cicd/app/logger"
	"go-cicd/app/restful/base"
	"go-cicd/app/utils"
	"net/http"
)

// HandleLogin create token and send to client
func HandleLogin(user *entity.User, w http.ResponseWriter) {
	tokenGenerator := authenticate.ResolveTokenGenerator()
	tokenStr, err := tokenGenerator.GenTokenForUser(user)

	if err != nil {
		logger.Error("Gen token error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		session := utils.GenUUIDString()
		authenticateInfo := map[string]interface{}{
			"token":   tokenStr,
			"session": session,
		}
		loginBody := map[string]interface{}{
			"authenticate": authenticateInfo,
			"user":         model.NewUserFromEntity(user),
		}

		response := model.APIResponse{
			Body: loginBody,
		}

		base.ResponseToClient(&response, w)
	}
}
