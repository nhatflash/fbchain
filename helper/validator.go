package helper

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const PHONE_PATTERN string = "^(0|\\+84|84)?(3[2-9]|5[6|8|9]|7[0|6-9]|8[1-9]|9[0-9])[0-9]{7}$"
const IDENTITY_PATTERN string = "^0[0-9]{11}$"
const NAME_PATTERN string = "^[a-zA-ZàáảãạăằắẳẵặâầấẩẫậèéẻẽẹêềếểễệìíỉĩịòóỏõọôồốổỗộơờớởỡợùúủũụưừứửữựỳýỷỹỵđÀÁẢÃẠĂẰẮẲẴẶÂẦẤẨẪẬÈÉẺẼẸÊỀẾỂỄỆÌÍỈĨỊÒÓỎÕỌÔỒỐỔỖỘƠỜỚỞỠỢÙÚỦŨỤƯỪỨỬỮỰỲÝỶỸỴĐ\\s]+$"


var PhoneNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if ok {
		_, err := regexp.MatchString(PHONE_PATTERN, phone)
		if err != nil {
			return false
		}
	}
	return true
}

var IdentityNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	identity, ok := fl.Field().Interface().(string)
	if ok {
		_, err := regexp.MatchString(IDENTITY_PATTERN, identity)
		if err != nil {
			return false
		}
	}
	return true
}

var NameValidator validator.Func = func(fl validator.FieldLevel) bool {
	name, ok := fl.Field().Interface().(string)
	if ok {
		_, err := regexp.MatchString(NAME_PATTERN, name)
		if err != nil {
			return false
		}
	}
	return true
}

var PostalCodeValidator validator.Func = func(fl validator.FieldLevel) bool {
	postalCode, ok := fl.Field().Interface().(string)
	if ok {
		if len(postalCode) > 10 {
			return false
		}
	}
	return true
}