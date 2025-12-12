package helper

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	PhonePattern string = "^(0|\\+84|84)?(3[2-9]|5[6|8|9]|7[0|6-9]|8[1-9]|9[0-9])[0-9]{7}$"
	IdentityPattern string = "^0[0-9]{11}$"
	NamePattern string = "^[a-zA-ZàáảãạăằắẳẵặâầấẩẫậèéẻẽẹêềếểễệìíỉĩịòóỏõọôồốổỗộơờớởỡợùúủũụưừứửữựỳýỷỹỵđÀÁẢÃẠĂẰẮẲẴẶÂẦẤẨẪẬÈÉẺẼẸÊỀẾỂỄỆÌÍỈĨỊÒÓỎÕỌÔỒỐỔỖỘƠỜỚỞỠỢÙÚỦŨỤƯỪỨỬỮỰỲÝỶỸỴĐ\\s]+$"
	PricePattern = "^[0-9]+$"
)


var PhoneNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if ok {
		match, err := regexp.MatchString(PhonePattern, phone)
		if err != nil || !match {
			return false
		}
	}
	return true
}

var IdentityNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	identity, ok := fl.Field().Interface().(string)
	if ok {
		match, err := regexp.MatchString(IdentityPattern, identity)
		if err != nil || !match {
			return false
		}
	}
	return true
}

var NameValidator validator.Func = func(fl validator.FieldLevel) bool {
	name, ok := fl.Field().Interface().(string)
	if ok {
		match, err := regexp.MatchString(NamePattern, name)
		if err != nil || !match {
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

var PositiveNumberValidator validator.Func = func(fl validator.FieldLevel) bool {
	positiveNum, ok := fl.Field().Interface().(int)
	if ok {
		if positiveNum < 0 {
			return false
		}
	}
	return true
}


var PriceValidator validator.Func = func(fl validator.FieldLevel) bool {
	price, ok := fl.Field().Interface().(string)
	if ok {
		if !strings.Contains(price, ".") {
			return false
		}
		parts := strings.Split(price, ".")
		if len(parts) != 2 || len(parts[1]) != 2 {
			return false
		}
		mValue, vErr := regexp.MatchString(PricePattern, parts[0])
		mPrecision, pErr := regexp.MatchString(PricePattern, parts[1])
		if vErr != nil || pErr != nil {
			return false
		}
		if !mValue || !mPrecision {
			return false
		}
	}
	return true
}

