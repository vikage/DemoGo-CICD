package repository

import (
	"fmt"
	"go-cicd/app/database"
	"go-cicd/app/domain/entity"
	"go-cicd/app/utils"
	"go-cicd/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userRepoTestSuite struct {
	suite.Suite
	dbClient     database.Client
	dbClientMock *mocks.Client
	userRepo     UserRepo
}

func (suite *userRepoTestSuite) SetupTest() {
	database.ClearSimulatorData()
	suite.dbClient = database.NewSimulatorClient()
	suite.userRepo = NewUserRepo(suite.dbClient)
	suite.dbClientMock = new(mocks.Client)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(userRepoTestSuite))
}

// ===========Find user by email tests ==========
func (suite *userRepoTestSuite) TestFindUserByEmailSuccess() {
	fakeUser := entity.User{
		ID:    utils.GenUUIDString(),
		Email: "thanhdaihiep94@gmail.com",
	}

	err := suite.userRepo.AddUser(&fakeUser)
	assert.Nil(suite.T(), err)

	userFromDB, err := suite.userRepo.FindUserByEmail("thanhdaihiep94@gmail.com")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), fakeUser.ID, userFromDB.ID)
}

func (suite *userRepoTestSuite) TestFindUserByEmailError() {
	fakeUser := entity.User{
		ID:    utils.GenUUIDString(),
		Email: "thanhdaihiep94@gmail.com",
	}

	database := new(mocks.Database)
	collection := new(mocks.Collection)
	findOneResult := new(mocks.SingleResult)

	suite.dbClientMock.On("Database", mock.Anything).Return(database)
	database.On("Collection", mock.Anything).Return(collection)
	collection.On("FindOne", mock.Anything, mock.Anything).Return(findOneResult)
	findOneResult.On("Decode", mock.Anything).Return(fmt.Errorf("fake error"))

	suite.userRepo = NewUserRepo(suite.dbClientMock)

	userFromDB, err := suite.userRepo.FindUserByEmail(fakeUser.Email)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), userFromDB)

	suite.dbClientMock.AssertExpectations(suite.T())
	database.AssertExpectations(suite.T())
	collection.AssertExpectations(suite.T())
}

func (suite *userRepoTestSuite) TestFindUserByEmailNotFound() {
	fakeUser := entity.User{
		ID:    utils.GenUUIDString(),
		Email: "thanhdaihiep94@gmail.com",
	}

	userFromDB, err := suite.userRepo.FindUserByEmail(fakeUser.Email)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), userFromDB)
}
