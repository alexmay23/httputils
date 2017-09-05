package httputils

import (
"errors"
"fmt"
"github.com/johngb/langreg"
"gopkg.in/mgo.v2/bson"
"log"
"net/url"
"strings"
"time"
)

type Validator func(value interface{}) error

func NotEmptyValidator(key string) Validator {
	return func(value interface{}) error {
		if value == nil {
			return Error{key, "Field is required", "REQUIRED_FIELD_ERROR", nil}
		}
		return nil
	}
}

func StringValidator(key string) Validator {
	return func(value interface{}) error {
		_, ok := value.(string)
		if !ok {
			return Error{key, " Should be string", "TYPE_ERROR", []string{"string"}}
		}
		return nil
	}
}

func FloatValidator(key string) Validator {
	return func(value interface{}) error {
		_, ok := value.(float64)
		if !ok {
			return Error{key, " Should be float", "TYPE_ERROR", []string{"float"}}
		}
		return nil
	}
}

type FloatRange struct {
	upper  float64
	bottom float64
}

func FloatInRangeValidator(key string, floatRange FloatRange) Validator {
	return func(value interface{}) error {
		float := value.(float64)
		if float > floatRange.upper || float < floatRange.bottom {
			return Error{key, "Invalid float", "FLOAT_RANGE_ERROR", nil}
		}
		return nil
	}
}

func ObjectIDValidator(key string) Validator {
	return func(value interface{}) error {
		str := value.(string)
		if !bson.IsObjectIdHex(str) {
			return Error{key, " Should be object id", "TYPE_ERROR", []string{"ObjectId"}}
		}
		return nil
	}
}

func StringLengthValidator(length int, key string) Validator {

	return func(value interface{}) error {
		stringValue := value.(string)
		if len(stringValue) < length {
			return Error{key, fmt.Sprintf("%@ should be minimum %d characters", strings.ToUpper(key), length),
						 "STRING_LENGTH_ERROR", []string{key, "5"}}

		}
		return nil
	}
}

func StringArrayValidator(key string, each []Validator) Validator {
	return func(value interface{}) error {
		array, ok := value.([]string)
		if !ok {
			return Error{key, "Should be array", "TYPE_ERROR", []string{"array"}}
		}
		for _, item := range array {
			for _, validator := range each {
				err := validator(item)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func LanguageValidator(key string) Validator {
	return func(value interface{}) error {
		stringValue := value.(string)
		if !langreg.IsValidLanguageCode(stringValue) {
			return Error{key, "Invalid language", "INVALID_LANGUAGE_ERROR", []string{stringValue}}

		}
		return nil
	}
}

func URLValidator(key string) Validator {
	return func(value interface{}) error {
		stringValue := value.(string)

		_, err := url.Parse(stringValue)
		if err != nil {
			return Error{key, "Invalid url", "INVALID_URL_ERROR", nil}
		}
		return nil
	}
}

func SexValidator(key string) Validator {
	return StringContainsValidator(key, []string{"male", "female"})
}

func StringContainsValidator(key string, values []string) Validator {
	return func(value interface{}) error {
		stringValue := value.(string)
		contains := false
		for _, item := range values {
			if item == stringValue {
				contains = true
				break
			}
		}
		if !contains {
			return Error{key, fmt.Sprintf("Invalid %s", key),
						 fmt.Sprintf("INVALID_%s_ERROR", strings.ToUpper(key)), nil}
		}
		return nil
	}
}

func TimezoneValidator(key string) Validator {
	return func(value interface{}) error {
		stringValue := value.(string)
		_, err := time.LoadLocation(stringValue)
		if err != nil {
			return Error{key, "Invalid timezone", "INVALID_TIMEZONE_ERROR", nil}
		}
		return nil
	}
}

func DateTimeValidator(key string, t *time.Time) Validator {
	return func(value interface{}) error {
		var err error
		switch value.(type) {
		case string:
			*t, err = time.Parse(time.RFC3339, value.(string))
		case float64:
			log.Print("SOME!", value)
			*t = time.Unix(int64(value.(float64)), 0)
		default:
			err = errors.New("Invalid datetime")
		}

		if err != nil {
			return Error{key, "Invalid datetime", "INVALID_DATETIME_ERROR", nil}
		}
		return nil
	}
}

func CountryValidator(key string) Validator {
	return func(value interface{}) error {
		stringValue := value.(string)
		if !langreg.IsValidRegionCode(stringValue) {
			return Error{key, "Invalid country", "INVALID_COUNTRY_ERROR", nil}
		}
		return nil
	}
}

func RequiredStringValidators(key string, validators ...Validator) []Validator {
	arr := []Validator{NotEmptyValidator(key), StringValidator(key)}
	return append(arr, validators...)
}

func RequiredFloatValidators(key string, validators ...Validator) []Validator {
	arr := []Validator{NotEmptyValidator(key), FloatValidator(key)}
	return append(arr, validators...)
}

func ValidateValue(value interface{}, validators []Validator) []Error {
	errs := []Error{}
	for _, validator := range validators {
		err := validator(value)
		if err != nil {
			errs = append(errs, err.(Error))
			break
		}
	}
	return errs
}

type VMap map[string][]Validator

func ValidateMap(dictionary map[string]interface{}, validatorMap VMap) []Error {
	errs := []Error{}
	for key, validators := range validatorMap {
		errs = append(errs, ValidateValue(dictionary[key], validators)...)
	}
	return errs
}
