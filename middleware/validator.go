package logiamdw

import (
	"github.com/go-playground/validator/v10"
	logiapkg "github.com/yusologia/go-core/v2/pkg"
	logiares "github.com/yusologia/go-core/v2/response"
	"strings"
	"time"
)

type Validator struct{}

func (v Validator) Make(rules interface{}) {
	err := logiapkg.LogiaValidate.Struct(rules)
	if err != nil {
		var attributes []interface{}
		for _, e := range err.(validator.ValidationErrors) {
			attributes = append(attributes, map[string]interface{}{
				"param":   e.Field(),
				"message": getMessage(e.Error()),
			})
		}

		logiares.ErrLogiaValidation(attributes)
	}
}

func (v Validator) RegisterValidation(callback func(validate *validator.Validate)) {
	logiapkg.LogiaValidate = validator.New()

	_ = logiapkg.LogiaValidate.RegisterValidation("date_ddmmyyyy", dateDDMMYYYYValidation)
	_ = logiapkg.LogiaValidate.RegisterValidation("time_hhmm", dateHHMMValidation)
	_ = logiapkg.LogiaValidate.RegisterValidation("time_hhmmss", dateHHMMSSValidation)

	callback(logiapkg.LogiaValidate)
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
