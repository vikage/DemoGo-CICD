package model

// APIResponse define fields of API response structure
type APIResponse struct {
	ErrorCode int         `json:"error_code"`
	Body      interface{} `json:"body,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// APIError define error field for API
type APIError int

const (
	// APIErrorSuccess success, no error
	APIErrorSuccess = 0
	// APIErrorWrongUserOrPassword wrong user or password
	APIErrorWrongUserOrPassword = 3
	// APIErrorUserNotFound user not found in system
	APIErrorUserNotFound = 4
	// APIErrorNeedChangePassword user logged in need change password
	APIErrorNeedChangePassword = 5
	// APIErrorResetPasswordRequestExpire request reset passwod did expire
	APIErrorResetPasswordRequestExpire = 6
	// APIErrorResetPasswordNotMatch not match password in db
	APIErrorResetPasswordNotMatch = 7
	// APIErrorLoginThirdpartyTokenExpire thirdparty token expire
	APIErrorLoginThirdpartyTokenExpire = 8
	// APIErrorUserExist user did exist in database
	APIErrorUserExist = 9
	// APIErrorPasswordNotMatch password not match with db
	APIErrorPasswordNotMatch = 10
	// APIErrorReferenceNotMatch reference code not match
	APIErrorReferenceNotMatch = 11
	// APIErrorDuplicate record already exist
	APIErrorDuplicate = 12
	// APIErrorOutOfStorage out of storage. Need upgrade storage
	APIErrorOutOfStorage = 13
	// APIErrorAccountNotAllow not allow doing some operation. Ex account login google can not reset password
	APIErrorAccountNotAllow = 14
	// APIErrorDidJoinAlbumBefore user did join before
	APIErrorDidJoinAlbumBefore = 15

	// APIErrorCanNotProcessResource can not process resource user upload like image/video
	APIErrorCanNotProcessResource = 201

	// APIErrorBadRequest bad request
	APIErrorBadRequest = 400
	// APIErrorForbidden forbidden
	APIErrorForbidden = 401
	// APIErrorNotFound not found resource
	APIErrorNotFound = 404
)
