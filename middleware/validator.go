package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/yusologia/go-core/config"
	"github.com/yusologia/go-core/response/error"
	"net/http"
	"strings"
	"time"
)

type Validator struct{}

func (v Validator) Make(r *http.Request, rules interface{}) {
	err := config.LogiaValidate.Struct(rules)
	if err != nil {
		var attributes []interface{}
		for _, e := range err.(validator.ValidationErrors) {
			attributes = append(attributes, map[string]interface{}{
				"param":   e.Field(),
				"message": getMessage(e.Error()),
			})
		}

		error.ErrLogiaValidation(attributes)
	}
}

func (v Validator) RegisterValidation(callback func(validate *validator.Validate)) {
	config.LogiaValidate = validator.New()

	_ = config.LogiaValidate.RegisterValidation("date_ddmmyyyy", dateDDMMYYYYValidation)
	_ = config.LogiaValidate.RegisterValidation("time_hhmm", dateHHMMValidation)
	_ = config.LogiaValidate.RegisterValidation("time_hhmmss", dateHHMMSSValidation)

	callback(config.LogiaValidate)
}

func getMessage(errMsg string) string {
	splitMsg := strings.Split(errMsg, ":")
	key := 0
	if len(splitMsg) == 3 {
		key = 2
	} else if len(splitMsg) == 2 {
		key = 1
	}

	return splitMsg[key]
}

func dateDDMMYYYYValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	_, err := time.Parse("02/01/2006", field)
	return err == nil
}

func dateHHMMValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	_, err := time.Parse("15:04", field)
	return err == nil
}

func dateHHMMSSValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	_, err := time.Parse("15:04:05", field)
	return err == nil
}
