package authenticate

import (
	"go-cicd/app/di"
	"go-cicd/app/di/gdi"
	"reflect"
)

var (
	// TokenGeneratorType type for authenticate.TokenGenerator
	TokenGeneratorType = reflect.TypeOf((*TokenGenerator)(nil)).Elem()
	// TokenDecoderType type for authenticate.TokenDecoderType
	TokenDecoderType = reflect.TypeOf((*TokenDecoder)(nil)).Elem()
)

// ResolveTokenGenerator get token generator from di
func ResolveTokenGenerator() TokenGenerator {
	generator, err := di.DefaultContainer.Resolve(TokenGeneratorType)
	if err != nil {
		return nil
	}

	return generator.(TokenGenerator)
}

// ResolveTokenDecoder get token decoder from di
func ResolveTokenDecoder() TokenDecoder {
	decoder, err := di.DefaultContainer.Resolve(TokenDecoderType)
	if err != nil {
		return nil
	}

	return decoder.(TokenDecoder)
}

// RegisterDependencyInContainer register database dependency
func RegisterDependencyInContainer(container *gdi.Container) {
	container.Register(TokenGeneratorType, NewTokenGenerator)
	container.Register(TokenDecoderType, NewTokenDecoder)
}
