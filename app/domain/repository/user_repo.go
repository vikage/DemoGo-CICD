package repository

import (
	"context"
	"go-cicd/app/database"
	"go-cicd/app/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepo define user repository
type UserRepo interface {
	SetSessionContext(sessionConext database.SessionContext)
	AddUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
	FindUserByID(id string) (*entity.User, error)
}

// UserRepoImpl implement for user repo
type UserRepoImpl struct {
	client         database.Client
	sessionContext database.SessionContext
}

// NewUserRepo Create repo
func NewUserRepo(dbClient database.Client) UserRepo {
	return &UserRepoImpl{client: dbClient}
}

// SetSessionContext set session context
func (repo *UserRepoImpl) SetSessionContext(sessionContext database.SessionContext) {
	repo.sessionContext = sessionContext
}

// AddUser add user to database
func (repo *UserRepoImpl) AddUser(user *entity.User) error {
	collection := repo.getCollection()

	ctx, cancel := repo.getContext()
	defer cancel()
	_, err := collection.InsertOne(ctx, user)

	return err
}

// FindUserByID query user in database by id
func (repo *UserRepoImpl) FindUserByID(id string) (*entity.User, error) {
	filter := bson.M{"_id": id}
	user, err := repo.queryUserByFilter(filter)

	return user, err
}

// FindUserByEmail find user by email in database
func (repo *UserRepoImpl) FindUserByEmail(email string) (*entity.User, error) {
	filter := bson.M{"email": email}
	user, err := repo.queryUserByFilter(filter)

	return user, err
}

func (repo *UserRepoImpl) queryUserByFilter(filter bson.M) (*entity.User, error) {
	collection := repo.getCollection()
	ctx, cancel := repo.getContext()
	defer cancel()

	var user entity.User
	result := collection.FindOne(ctx, filter)
	err := result.Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (repo *UserRepoImpl) getCollection() database.Collection {
	return repo.client.Database("go-cicd").Collection("User")
}

func (repo *UserRepoImpl) getContext() (context.Context, context.CancelFunc) {
	if repo.sessionContext != nil {
		return repo.sessionContext.Context(), func() {}
	}

	return database.MongoTimeoutContext()
}
