package authenticate

import (
	"fmt"
	"go-cicd/app/di"
	"go-cicd/app/domain/entity"
	"go-cicd/app/domain/repository"
	"go-cicd/app/utils"
	"go-cicd/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type tokenTestSuite struct {
	suite.Suite
	userRepoMock *mocks.UserRepo
}

func (suite *tokenTestSuite) SetupTest() {
	suite.userRepoMock = new(mocks.UserRepo)
	di.DefaultContainer.Register(repository.UserRepoType, func() repository.UserRepo {
		return suite.userRepoMock
	})
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(tokenTestSuite))
}

func (suite *tokenTestSuite) TestGenTokenAndPassSuccess() {
	user := entity.User{
		ID: utils.GenUUIDString(),
	}

	suite.userRepoMock.On("FindUserByID", mock.Anything).Return(&user, nil)
	generator := NewTokenGenerator()
	decoder := NewTokenDecoder()

	token, err := generator.GenTokenForUser(&user)
	if err != nil {
		suite.T().Fail()
	}

	userFromToken, err := decoder.UserFromToken(token)
	if err != nil {
		suite.T().Fail()
	}

	assert.Equal(suite.T(), user.ID, userFromToken.ID)
}

func (suite *tokenTestSuite) TestCaseUserRepoNotFoundUser() {
	user := entity.User{
		ID: utils.GenUUIDString(),
	}

	suite.userRepoMock.On("FindUserByID", mock.Anything).Return(nil, nil)
	generator := NewTokenGenerator()
	decoder := NewTokenDecoder()

	token, err := generator.GenTokenForUser(&user)
	if err != nil {
		suite.T().Fail()
	}

	_, err = decoder.UserFromToken(token)
	assert.NotNil(suite.T(), err)
}

func (suite *tokenTestSuite) TestCaseUserRepoReturnError() {
	user := entity.User{
		ID: utils.GenUUIDString(),
	}

	suite.userRepoMock.On("FindUserByID", mock.Anything).Return(nil, fmt.Errorf("Mock error"))
	generator := NewTokenGenerator()
	decoder := NewTokenDecoder()

	token, err := generator.GenTokenForUser(&user)
	if err != nil {
		suite.T().Fail()
	}

	_, err = decoder.UserFromToken(token)
	assert.NotNil(suite.T(), err)
}

func (suite *tokenTestSuite) TestCaseAccountKeyNotMatch() {
	user := entity.User{
		ID:         utils.GenUUIDString(),
		AccountKey: utils.GenUUIDString(),
	}

	userChangeKey := user
	userChangeKey.AccountKey = "Changed"

	suite.userRepoMock.On("FindUserByID", mock.Anything).Return(&userChangeKey, nil)
	generator := NewTokenGenerator()
	decoder := NewTokenDecoder()

	token, err := generator.GenTokenForUser(&user)
	if err != nil {
		suite.T().Fail()
	}

	_, err = decoder.UserFromToken(token)
	assert.NotNil(suite.T(), err)
}

func (suite *tokenTestSuite) TestCaseParseTokenInvalid() {
	decoder := NewTokenDecoder()
	_, err := decoder.UserFromToken("Invalid token")
	assert.NotNil(suite.T(), err)
}

func (suite *tokenTestSuite) TestParseTokenExpire() {
	user := entity.User{
		ID: utils.GenUUIDString(),
	}

	suite.userRepoMock.On("FindUserByID", mock.Anything).Return(&user, nil)
	generator := NewTokenGenerator()
	decoder := NewTokenDecoder()

	token, err := generator.GenTokenForUserWithExpireTime(&user, time.Now().Add(-time.Second))
	if err != nil {
		suite.T().Fail()
	}

	_, err = decoder.UserFromToken(token)
	assert.NotNil(suite.T(), err)
}
