package httputils

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/ti/mdb"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		params := make(map[string]string)
		for _, value := range ps {
			params[value.Key] = value.Value
		}
		h.ServeHTTP(w, SetInContext(params, "params", r))
	}
}

func SetInContext(value interface{}, key interface{}, req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

func ConvertMapToValue(value interface{}, jsonMap map[string]interface{}) error {
	data, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func JSON(w http.ResponseWriter, value interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

func DefaultMiddlewaresFactory(secret string) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		return AccessMiddlewareFactory(secret)(RecoverMiddleware(LoggingMiddleware(next)))
	}
	return f
}

func UnwrapOrDefault(value *int, d int) int {
	if value != nil {
		return *value
	}
	return d
}


func AccessMiddlewareFactory(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Secret") != secret {
				HTTP403().Write(w)
				return
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func WriteResponseOrError(w http.ResponseWriter, code int, response interface{}, err error) {
	if err != nil {
		err.(ServerError).Write(w)
		return
	}
	JSON(w, response, code)
}

func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				raise500(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Now() int64 {
	return int64(time.Now().Unix())
}

func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func GetBody(req *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(req.Body)
	var _map map[string]interface{}
	err := decoder.Decode(&_map)
	defer req.Body.Close()
	if err != nil {
		return nil, HTTP400()
	}
	return _map, nil
}

func ValidateBody(body map[string]interface{}, validatorMap VMap) (map[string]interface{}, error) {
	errs := ValidateMap(body, validatorMap)
	if len(errs) > 0 {
		return nil, ServerError{400, Errors{Errors: errs}}
	}
	return body, nil
}

func GetValidatedBody(req *http.Request, validatorMap VMap) (map[string]interface{}, error) {
	body, err := GetBody(req)
	if err != nil {
		return nil, err
	}
	return ValidateBody(body, validatorMap)
}

func MapKeys(m VMap) []string {
	keys := []string{};
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func GetValidatedURLParameters(req *http.Request, validatorMap VMap) (map[string]interface{}, error) {
	reqValues := make(map[string]interface{})
	for _, key := range MapKeys(validatorMap) {
		value := GetValueFromURLInRequest(req, key)
		if value == nil {
			reqValues[key] = nil
		} else {
			reqValues[key] = *value;
		}
	}
	errs := ValidateMap(reqValues, validatorMap);
	if len(errs) > 0 {
		return nil, ServerError{400, Errors{Errors: errs}}
	}
	return reqValues, nil
}

func ApplySkipLimit(query *mdb.Query, skip *int, limit *int) *mdb.Query {
	if skip != nil {
		query.Skip(*skip)
	}
	if limit != nil {
		query.Limit(*limit)
	}
	return query
}

func GetValueFromURLInRequest(r *http.Request, key string) *string {
	params := r.Context().Value("params").(map[string]string)
	var value string
	if len(params) > 0 {
		value = params[key]
	}
	if len(value) == 0 {
		value = r.URL.Query().Get(key)
	}

	if len(value) == 0 {
		return nil
	}
	return &value
}

func GetObjectIdFromURLInRequest(r *http.Request, key string) *bson.ObjectId {
	id := GetValueFromURLInRequest(r, key)
	if id == nil {
		return nil
	}
	if !bson.IsObjectIdHex(*id) {
		return nil
	}
	objectID := bson.ObjectIdHex(*id)
	return &objectID
}

func contains(array []string, element string) bool {
	for _, value := range array {
		if value == element {
			return true
		}
	}
	return false
}

func Find(collection *mdb.Collection, q bson.M, skip *int, limit *int) (*interface{}, int) {
	results := new(interface{})
	query := collection.Find(q)
	count, err := query.Count()
	if err != nil {
		panic(err)
	}
	err = ApplySkipLimit(query, skip, limit).All(&results)
	if err != nil {
		panic(err)
	}
	return results, count
}

func IntParameterFromRequest(r *http.Request, name string) *int {
	s := GetValueFromURLInRequest(r, name)
	if s == nil {
		return nil
	}
	i64, err := strconv.ParseInt(*s, 10, 64)
	if err != nil {
		return nil
	}
	i := int(i64)
	return &i
}

func FloatParameterFromRequest(r *http.Request, name string) *float64 {
	s := GetValueFromURLInRequest(r, name)
	if s == nil {
		return nil
	}
	float, err := strconv.ParseFloat(*s, 64);
	if err != nil {
		return nil
	}
	return &float
}
