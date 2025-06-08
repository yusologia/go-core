package logiares

import (
	"net/http"
)

func Error(code int, message string, internalMsg string, bug bool, attributes any) {
	panic(&ResponseError{
		Status: StatusError{
			Bug: bug,
			Status: Status{
				Code:        code,
				Message:     message,
				InternalMsg: internalMsg,
				Attributes:  attributes,
			},
		},
	})
}

func ErrLogiaUnauthenticated(internalMsg string) {
	Error(http.StatusUnauthorized, "Unauthenticated.", internalMsg, false, nil)
}

func ErrLogiaBadRequest(internalMsg string) {
	Error(http.StatusBadRequest, "Bad request!", internalMsg, false, nil)
}

func ErrLogiaPayloadVeryLarge(internalMsg string) {
	Error(http.StatusRequestEntityTooLarge, "Your payload very large!", internalMsg, false, nil)
}

func ErrLogiaValidation(attributes []interface{}) {
	Error(http.StatusBadRequest, "Missing Required Parameter", "", false, attributes)
}

func ErrLogiaNotFound(internalMsg string) {
	Error(http.StatusNotFound, "Data not found", internalMsg, false, nil)
}

func ErrLogiaUploadFile(internalMsg string) {
	Error(http.StatusInternalServerError, "Unable to upload file", internalMsg, false, nil)
}

func ErrLogiaDeleteFile(internalMsg string) {
	Error(http.StatusInternalServerError, "Unable to delete file", internalMsg, false, nil)
}

func ErrLogiaUUID(internalMsg string) {
	Error(http.StatusInternalServerError, "Unable to generate uuid", internalMsg, false, nil)
}

func ErrLogiaAPI(internalMsg string) {
	Error(http.StatusInternalServerError, "Calling external api is invalid!", internalMsg, false, nil)
}
