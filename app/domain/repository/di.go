package repository

import (
	"go-cicd/app/database"
	"go-cicd/app/di"
	"go-cicd/app/di/gdi"
	"reflect"
)

var (
	// UserRepoType use for UserRepo
	UserRepoType = reflect.TypeOf((*UserRepo)(nil)).Elem()
	// AlbumRepoType use for AlbumRepo
)

// RegisterDependencyInContainer register dependency in DI container
func RegisterDependencyInContainer(container *gdi.Container) {
	container.Register(UserRepoType, NewUserRepo)
}

// ResolveUserRepo resolve user repo
func ResolveUserRepo(client database.Client) UserRepo {
	repo, err := di.DefaultContainer.Resolve(UserRepoType, client)
	if err != nil {
		panic(err)
	}

	return repo.(UserRepo)
}
