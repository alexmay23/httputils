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


func BoolValidator(key string) Validator {
	return func(value interface{}) error {
		_, ok := value.(bool)
		if !ok {
			return Error{key, " Should be bool", "TYPE_ERROR", []string{"bool"}}
		}
		return nil
	}
}

func IntValidator(key string) Validator {
	return func(value interface{}) error {
		_, ok := value.(int)
		_, ok = value.(int64)
		if !ok {
			return Error{key, " Should be int", "TYPE_ERROR", []string{"int"}}
		}
		return nil
	}
}


type FloatRange struct {
	Upper  *float64
	Bottom *float64
}

type IntRange struct {
	Upper  *int
	Bottom *int
}

type Int64Range struct {
	Upper  *int64
	Bottom *int64
}

func FloatInRangeValidator(key string, floatRange FloatRange) Validator {
	return func(value interface{}) error {
		float := value.(float64)
		err := Error{key, "Invalid float", "FLOAT_RANGE_ERROR", nil}
		if floatRange.Upper != nil && *floatRange.Upper < float {
			return err
		}
		if floatRange.Bottom != nil && *floatRange.Bottom > float {
			return err
		}
		return nil
	}
}

func IntInRangeValidator(key string, intRange IntRange) Validator {
	return func(value interface{}) error {
		intValue := value.(int)
		err := Error{key, "Invalid int", "INT_RANGE_ERROR", nil}
		if intRange.Upper != nil && *intRange.Upper < intValue {
			return err
		}
		if intRange.Bottom != nil && *intRange.Bottom > intValue {
			return err
		}
		return nil
	}
}


func Int64InRangeValidator(key string, intRange Int64Range) Validator {
	return func(value interface{}) error {
		intValue := value.(int64)
		err := Error{key, "Invalid int", "INT_RANGE_ERROR", nil}
		if intRange.Upper != nil && *intRange.Upper < intValue {
			return err
		}
		if intRange.Bottom != nil && *intRange.Bottom > intValue {
			return err
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

func ArrayValidator(key string) Validator {
	return func(value interface{}) error {
		_, ok := value.([]interface{})
		if !ok {
			return Error{key, "Should be array", "TYPE_ERROR", []string{"array"}}
		}
		return nil
	}
}

func StringArrayValidator(key string, each []Validator) Validator {
	return func(value interface{}) error {
		values := value.([]interface{})
		strArr := []string{}
		for _, item := range values {
			str, ok := item.(string)
			if !ok {
				return Error{key, "Should be string in array", "TYPE_ERROR", []string{"string", "array"}}
			}
			strArr = append(strArr, str)
		}
		for _, item := range strArr {
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

var timezones = []string{"Africa/Abidjan", "Africa/Accra", "Africa/Addis_Ababa", "Africa/Algiers", "Africa/Asmara",
	"Africa/Bamako", "Africa/Bangui", "Africa/Banjul", "Africa/Bissau", "Africa/Blantyre", "Africa/Brazzaville",
	"Africa/Bujumbura", "Africa/Cairo", "Africa/Casablanca", "Africa/Ceuta", "Africa/Conakry", "Africa/Dakar",
	"Africa/Dar_es_Salaam", "Africa/Djibouti", "Africa/Douala", "Africa/El_Aaiun", "Africa/Freetown", "Africa/Gaborone",
	"Africa/Harare", "Africa/Johannesburg", "Africa/Juba", "Africa/Kampala", "Africa/Khartoum", "Africa/Kigali",
	"Africa/Kinshasa", "Africa/Lagos", "Africa/Libreville", "Africa/Lome", "Africa/Luanda", "Africa/Lubumbashi",
	"Africa/Lusaka", "Africa/Malabo", "Africa/Maputo", "Africa/Maseru", "Africa/Mbabane", "Africa/Mogadishu",
	"Africa/Monrovia", "Africa/Nairobi", "Africa/Ndjamena", "Africa/Niamey", "Africa/Nouakchott", "Africa/Ouagadougou",
	"Africa/Porto-Novo", "Africa/Sao_Tome", "Africa/Tripoli", "Africa/Tunis", "Africa/Windhoek", "America/Adak",
	"America/Anchorage", "America/Anguilla", "America/Antigua", "America/Araguaina", "America/Argentina/Buenos_Aires",
	"America/Argentina/Catamarca", "America/Argentina/Cordoba", "America/Argentina/Jujuy", "America/Argentina/La_Rioja",
	"America/Argentina/Mendoza", "America/Argentina/Rio_Gallegos", "America/Argentina/Salta", "America/Argentina/San_Juan",
	"America/Argentina/San_Luis", "America/Argentina/Tucuman", "America/Argentina/Ushuaia", "America/Aruba", "America/Asuncion",
	"America/Atikokan", "America/Bahia", "America/Bahia_Banderas", "America/Barbados", "America/Belem", "America/Belize",
	"America/Blanc-Sablon", "America/Boa_Vista", "America/Bogota", "America/Boise", "America/Cambridge_Bay",
	"America/Campo_Grande", "America/Cancun", "America/Caracas", "America/Cayenne", "America/Cayman",
	"America/Chicago", "America/Chihuahua", "America/Costa_Rica", "America/Creston", "America/Cuiaba",
	"America/Curacao", "America/Danmarkshavn", "America/Dawson", "America/Dawson_Creek", "America/Denver",
	"America/Detroit", "America/Dominica", "America/Edmonton", "America/Eirunepe", "America/El_Salvador",
	"America/Fort_Nelson", "America/Fortaleza", "America/Glace_Bay", "America/Godthab", "America/Goose_Bay",
	"America/Grand_Turk", "America/Grenada", "America/Guadeloupe", "America/Guatemala", "America/Guayaquil", "America/Guyana", "America/Halifax", "America/Havana", "America/Hermosillo", "America/Indiana/Indianapolis", "America/Indiana/Knox", "America/Indiana/Marengo", "America/Indiana/Petersburg", "America/Indiana/Tell_City", "America/Indiana/Vevay", "America/Indiana/Vincennes", "America/Indiana/Winamac", "America/Inuvik", "America/Iqaluit", "America/Jamaica", "America/Juneau", "America/Kentucky/Louisville", "America/Kentucky/Monticello", "America/Kralendijk", "America/La_Paz", "America/Lima", "America/Los_Angeles", "America/Lower_Princes", "America/Maceio", "America/Managua", "America/Manaus", "America/Marigot", "America/Martinique", "America/Matamoros", "America/Mazatlan", "America/Menominee", "America/Merida", "America/Metlakatla", "America/Mexico_City", "America/Miquelon", "America/Moncton", "America/Monterrey", "America/Montevideo", "America/Montreal", "America/Montserrat", "America/Nassau", "America/New_York", "America/Nipigon", "America/Nome", "America/Noronha", "America/North_Dakota/Beulah", "America/North_Dakota/Center", "America/North_Dakota/New_Salem", "America/Ojinaga", "America/Panama", "America/Pangnirtung", "America/Paramaribo", "America/Phoenix", "America/Port-au-Prince", "America/Port_of_Spain", "America/Porto_Velho", "America/Puerto_Rico", "America/Punta_Arenas", "America/Rainy_River", "America/Rankin_Inlet", "America/Recife", "America/Regina", "America/Resolute", "America/Rio_Branco", "America/Santa_Isabel", "America/Santarem", "America/Santiago", "America/Santo_Domingo", "America/Sao_Paulo", "America/Scoresbysund", "America/Shiprock", "America/Sitka", "America/St_Barthelemy", "America/St_Johns", "America/St_Kitts", "America/St_Lucia", "America/St_Thomas", "America/St_Vincent", "America/Swift_Current", "America/Tegucigalpa", "America/Thule", "America/Thunder_Bay", "America/Tijuana", "America/Toronto", "America/Tortola", "America/Vancouver", "America/Whitehorse", "America/Winnipeg", "America/Yakutat", "America/Yellowknife", "Antarctica/Casey", "Antarctica/Davis", "Antarctica/DumontDUrville", "Antarctica/Macquarie", "Antarctica/Mawson", "Antarctica/McMurdo", "Antarctica/Palmer", "Antarctica/Rothera", "Antarctica/South_Pole", "Antarctica/Syowa", "Antarctica/Troll", "Antarctica/Vostok", "Arctic/Longyearbyen", "Asia/Aden", "Asia/Almaty", "Asia/Amman", "Asia/Anadyr", "Asia/Aqtau", "Asia/Aqtobe", "Asia/Ashgabat", "Asia/Atyrau", "Asia/Baghdad", "Asia/Bahrain", "Asia/Baku", "Asia/Bangkok", "Asia/Barnaul", "Asia/Beirut", "Asia/Bishkek", "Asia/Brunei", "Asia/Calcutta", "Asia/Chita", "Asia/Choibalsan", "Asia/Chongqing", "Asia/Colombo", "Asia/Damascus", "Asia/Dhaka", "Asia/Dili", "Asia/Dubai", "Asia/Dushanbe", "Asia/Famagusta", "Asia/Gaza", "Asia/Harbin", "Asia/Hebron", "Asia/Ho_Chi_Minh", "Asia/Hong_Kong", "Asia/Hovd", "Asia/Irkutsk", "Asia/Jakarta", "Asia/Jayapura", "Asia/Jerusalem", "Asia/Kabul", "Asia/Kamchatka", "Asia/Karachi", "Asia/Kashgar", "Asia/Kathmandu", "Asia/Katmandu", "Asia/Khandyga", "Asia/Krasnoyarsk", "Asia/Kuala_Lumpur", "Asia/Kuching", "Asia/Kuwait", "Asia/Macau", "Asia/Magadan", "Asia/Makassar", "Asia/Manila", "Asia/Muscat", "Asia/Nicosia", "Asia/Novokuznetsk", "Asia/Novosibirsk", "Asia/Omsk", "Asia/Oral", "Asia/Phnom_Penh", "Asia/Pontianak", "Asia/Pyongyang", "Asia/Qatar", "Asia/Qyzylorda", "Asia/Rangoon", "Asia/Riyadh", "Asia/Sakhalin", "Asia/Samarkand", "Asia/Seoul", "Asia/Shanghai", "Asia/Singapore", "Asia/Srednekolymsk", "Asia/Taipei", "Asia/Tashkent", "Asia/Tbilisi", "Asia/Tehran", "Asia/Thimphu", "Asia/Tokyo", "Asia/Tomsk", "Asia/Ulaanbaatar", "Asia/Urumqi", "Asia/Ust-Nera", "Asia/Vientiane", "Asia/Vladivostok", "Asia/Yakutsk", "Asia/Yangon", "Asia/Yekaterinburg", "Asia/Yerevan", "Atlantic/Azores", "Atlantic/Bermuda", "Atlantic/Canary", "Atlantic/Cape_Verde", "Atlantic/Faroe", "Atlantic/Madeira", "Atlantic/Reykjavik", "Atlantic/South_Georgia", "Atlantic/St_Helena", "Atlantic/Stanley", "Australia/Adelaide", "Australia/Brisbane", "Australia/Broken_Hill", "Australia/Currie", "Australia/Darwin", "Australia/Eucla", "Australia/Hobart", "Australia/Lindeman", "Australia/Lord_Howe", "Australia/Melbourne", "Australia/Perth", "Australia/Sydney", "Europe/Amsterdam", "Europe/Andorra", "Europe/Astrakhan", "Europe/Athens", "Europe/Belgrade", "Europe/Berlin", "Europe/Bratislava", "Europe/Brussels", "Europe/Bucharest", "Europe/Budapest", "Europe/Busingen", "Europe/Chisinau", "Europe/Copenhagen", "Europe/Dublin", "Europe/Gibraltar", "Europe/Guernsey", "Europe/Helsinki", "Europe/Isle_of_Man", "Europe/Istanbul", "Europe/Jersey", "Europe/Kaliningrad", "Europe/Kiev", "Europe/Kirov", "Europe/Lisbon", "Europe/Ljubljana", "Europe/London", "Europe/Luxembourg", "Europe/Madrid", "Europe/Malta", "Europe/Mariehamn", "Europe/Minsk", "Europe/Monaco", "Europe/Moscow", "Europe/Oslo", "Europe/Paris", "Europe/Podgorica", "Europe/Prague", "Europe/Riga", "Europe/Rome", "Europe/Samara", "Europe/San_Marino", "Europe/Sarajevo", "Europe/Saratov", "Europe/Simferopol", "Europe/Skopje", "Europe/Sofia", "Europe/Stockholm", "Europe/Tallinn", "Europe/Tirane", "Europe/Ulyanovsk", "Europe/Uzhgorod", "Europe/Vaduz", "Europe/Vatican", "Europe/Vienna", "Europe/Vilnius", "Europe/Volgograd", "Europe/Warsaw", "Europe/Zagreb", "Europe/Zaporozhye", "Europe/Zurich", "GMT", "Indian/Antananarivo", "Indian/Chagos", "Indian/Christmas", "Indian/Cocos", "Indian/Comoro", "Indian/Kerguelen", "Indian/Mahe", "Indian/Maldives", "Indian/Mauritius", "Indian/Mayotte", "Indian/Reunion", "Pacific/Apia", "Pacific/Auckland", "Pacific/Bougainville", "Pacific/Chatham", "Pacific/Chuuk", "Pacific/Easter", "Pacific/Efate", "Pacific/Enderbury", "Pacific/Fakaofo", "Pacific/Fiji", "Pacific/Funafuti", "Pacific/Galapagos", "Pacific/Gambier", "Pacific/Guadalcanal", "Pacific/Guam", "Pacific/Honolulu", "Pacific/Johnston", "Pacific/Kiritimati", "Pacific/Kosrae", "Pacific/Kwajalein", "Pacific/Majuro", "Pacific/Marquesas", "Pacific/Midway", "Pacific/Nauru", "Pacific/Niue", "Pacific/Norfolk", "Pacific/Noumea", "Pacific/Pago_Pago", "Pacific/Palau", "Pacific/Pitcairn", "Pacific/Pohnpei", "Pacific/Ponape", "Pacific/Port_Moresby", "Pacific/Rarotonga", "Pacific/Saipan", "Pacific/Tahiti", "Pacific/Tarawa", "Pacific/Tongatapu", "Pacific/Truk", "Pacific/Wake", "Pacific/Wallis"}

func TimezoneValidator(key string) Validator {
	return func(value interface{}) error {
		stringValue := value.(string)
		if !contains(timezones, stringValue) {
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

func RequiredBoolValidators(key string, validators ...Validator) []Validator {
	arr := []Validator{NotEmptyValidator(key), BoolValidator(key)}
	return append(arr, validators...)
}

func RequiredIntValidators(key string, validators ...Validator) []Validator {
	arr := []Validator{NotEmptyValidator(key), IntValidator(key)}
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
