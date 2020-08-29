package utils

import (
	"encoding/json"
	"go-cicd/app/domain/model"
	"io"
)

// DecodeAPIResponse decode api response for test case verify
func DecodeAPIResponse(buffer io.Reader) (*model.APIResponse, error) {
	var response model.APIResponse
	err := json.NewDecoder(buffer).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
