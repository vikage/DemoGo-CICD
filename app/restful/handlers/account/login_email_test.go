package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-cicd/app/authenticate"
	"go-cicd/app/di"
	"go-cicd/app/domain/entity"
	"go-cicd/app/domain/model"
	"go-cicd/app/domain/repository"
	"go-cicd/app/utils"
	"go-cicd/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type loginEmailTestStuite struct {
	suite.Suite
	userRepoMock *mocks.UserRepo
}

func (suite *loginEmailTestStuite) SetupTest() {
	authenticate.RegisterDependencyInContainer(di.DefaultContainer)
	suite.userRepoMock = new(mocks.UserRepo)
	di.DefaultContainer.Register(repository.UserRepoType, func() repository.UserRepo {
		return suite.userRepoMock
	})
}

func TestLoginEmailTestSuite(t *testing.T) {
	suite.Run(t, new(loginEmailTestStuite))
}

func (suite *loginEmailTestStuite) TestLoginEmailSuccess() {
	fakeUser := entity.User{
		ID:         utils.GenUUIDString(),
		AccountKey: utils.GenUUIDString(),
		Email:      "email@example.com",
	}

	fakeUser.Password = utils.EncryptPassword("123456", fakeUser.ID)
	suite.userRepoMock.On("FindUserByEmail", mock.Anything).Return(&fakeUser, nil)

	responseRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"email":    "email@example.com",
		"password": "123456",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		suite.T().Fail()
		return
	}

	request, err := http.NewRequest("POST", "", bytes.NewBuffer(data))
	if err != nil {
		suite.T().Fail()
		return
	}

	LoginRequestHandler(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Result().StatusCode, "Status code must be 200")

	var response model.APIResponse
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	assert.Equal(suite.T(), model.APIErrorSuccess, response.ErrorCode, "Response error code must be success(0)")

	if bodyMap, ok := response.Body.(map[string]interface{}); ok {
		if !utils.ExistKeyInMap(bodyMap, []string{"authenticate", "user"}, []string{"map[string]interface {}", "map[string]interface {}"}) {
			suite.T().Error("Body must contain authenticate and user")
			suite.T().Fail()
			return
		}

		authenticateInfo := bodyMap["authenticate"].(map[string]interface{})
		if !utils.ExistKeyInMap(authenticateInfo, []string{"token", "session"}, []string{"string", "string"}) {
			suite.T().Error("Authenticate must contain token and session")
			suite.T().Fail()
			return
		}
	} else {
		suite.T().Error("Body must be a map[string]interface{}")
		suite.T().Fail()
	}
}

func (suite *loginEmailTestStuite) TestLoginEmailOrPasswordIncorrect() {
	fakeUser := entity.User{
		ID:         utils.GenUUIDString(),
		AccountKey: utils.GenUUIDString(),
		Email:      "email@example.com",
	}

	fakeUser.Password = utils.EncryptPassword("123456", fakeUser.ID)
	suite.userRepoMock.On("FindUserByEmail", mock.Anything).Return(&fakeUser, nil)

	responseRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"email":    "email@example.com",
		"password": "1234567",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		suite.T().Fail()
		return
	}

	request, err := http.NewRequest("POST", "", bytes.NewBuffer(data))
	if err != nil {
		suite.T().Fail()
		return
	}

	LoginRequestHandler(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Status code must be 401")

	var response model.APIResponse
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	assert.Equal(suite.T(), model.APIErrorWrongUserOrPassword, response.ErrorCode, "Response error code must be APIErrorWrongUserOrPassword")
}

func (suite *loginEmailTestStuite) TestLoginMissingEmail() {
	fakeUser := entity.User{
		ID:         utils.GenUUIDString(),
		AccountKey: utils.GenUUIDString(),
		Email:      "email@example.com",
	}

	fakeUser.Password = utils.EncryptPassword("123456", fakeUser.ID)
	suite.userRepoMock.On("FindUserByEmail", mock.Anything).Return(&fakeUser, nil)

	responseRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"password": "1234567",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		suite.T().Fail()
		return
	}

	request, err := http.NewRequest("POST", "", bytes.NewBuffer(data))
	if err != nil {
		suite.T().Fail()
		return
	}

	LoginRequestHandler(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Result().StatusCode, "Status code must be 400")
}

func (suite *loginEmailTestStuite) TestLoginMissingPasword() {
	fakeUser := entity.User{
		ID:         utils.GenUUIDString(),
		AccountKey: utils.GenUUIDString(),
		Email:      "email@example.com",
	}

	fakeUser.Password = utils.EncryptPassword("123456", fakeUser.ID)
	suite.userRepoMock.On("FindUserByEmail", mock.Anything).Return(&fakeUser, nil)

	responseRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"email": "email@example.com",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		suite.T().Fail()
		return
	}

	request, err := http.NewRequest("POST", "", bytes.NewBuffer(data))
	if err != nil {
		suite.T().Fail()
		return
	}

	LoginRequestHandler(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Result().StatusCode, "Status code must be 400")
}

func (suite *loginEmailTestStuite) TestLoginDBError() {
	suite.userRepoMock.On("FindUserByEmail", mock.Anything).Return(nil, fmt.Errorf("Fake error"))

	responseRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"email":    "email@example.com",
		"password": "123456",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		suite.T().Fail()
		return
	}

	request, err := http.NewRequest("POST", "", bytes.NewBuffer(data))
	if err != nil {
		suite.T().Fail()
		return
	}

	LoginRequestHandler(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Status code must be 500")
}

func (suite *loginEmailTestStuite) TestLoginUserNotFound() {
	suite.userRepoMock.On("FindUserByEmail", mock.Anything).Return(nil, nil)

	responseRecorder := httptest.NewRecorder()

	payload := map[string]interface{}{
		"email":    "email@example.com",
		"password": "1234567",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		suite.T().Fail()
		return
	}

	request, err := http.NewRequest("POST", "", bytes.NewBuffer(data))
	if err != nil {
		suite.T().Fail()
		return
	}

	LoginRequestHandler(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Status code must be 401")

	var response model.APIResponse
	json.NewDecoder(responseRecorder.Body).Decode(&response)
	assert.Equal(suite.T(), model.APIErrorWrongUserOrPassword, response.ErrorCode, "Response error code must be APIErrorWrongUserOrPassword")
}
