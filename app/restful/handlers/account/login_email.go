package account

import (
	"encoding/json"
	"go-cicd/app/database"
	"go-cicd/app/domain/model"
	"go-cicd/app/domain/repository"
	"go-cicd/app/logger"
	"go-cicd/app/restful/base"
	"go-cicd/app/restful/payloads"
	"go-cicd/app/utils"
	"net/http"
)

// LoginRequestHandler Handle login for Rest API
func LoginRequestHandler(w http.ResponseWriter, r *http.Request) {
	var payload payloads.LoginEmailPayload
	json.NewDecoder(r.Body).Decode(&payload)

	isValid, err := payload.Validate()
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		base.ResponseError(model.APIErrorBadRequest, err.Error(), w)
		logger.ErrorReq(r, "Validate err: %s", err)
		return
	}

	dbClient := database.ResolveDatabaseClient()
	userRepo := repository.ResolveUserRepo(dbClient)

	userEntity, err := userRepo.FindUserByEmail(payload.Email)

	if err != nil {
		logger.ErrorReq(r, "%s Access database error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userEntity != nil {
		if userEntity.Password == utils.EncryptPassword(payload.Password, userEntity.ID) {
			HandleLogin(userEntity, w)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	base.ResponseError(model.APIErrorWrongUserOrPassword, "Wrong email or password", w)
}
