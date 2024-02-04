package error

import (
	"github.com/yusologia/go-core/response"
	"net/http"
)

func ErrLogiaUnauthenticated(internalMsg string) {
	response.Error(http.StatusUnauthorized, "Unauthenticated.", internalMsg, nil)
}

func ErrLogiaBadRequest(internalMsg string) {
	response.Error(http.StatusBadRequest, "Bad request!", internalMsg, nil)
}

func ErrLogiaPayloadVeryLarge(internalMsg string) {
	response.Error(http.StatusRequestEntityTooLarge, "Your payload very large!", internalMsg, nil)
}

func ErrLogiaValidation(attributes []interface{}) {
	response.Error(http.StatusBadRequest, "Missing Required Parameter", "", attributes)
}

func ErrLogiaNotFound(internalMsg string) {
	response.Error(http.StatusNotFound, "Data not found", internalMsg, nil)
}
