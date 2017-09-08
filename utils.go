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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func JSON(w http.ResponseWriter, value interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json")
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

func DefaultMiddlewares(next http.Handler) http.Handler {
	return AccessMiddleware(RecoverMiddleware(LoggingMiddleware(next)))
}

func AccessMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Secret") != "Excellent" {
			HTTP403().Write(w)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				raise500(w)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Now()int{
	return int(time.Now().Unix())
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

func GetValidatedBody(req *http.Request, validatorMap VMap) (map[string]interface{}, error) {
	body, err := GetBody(req)
	if err != nil {
		return nil, err
	}
	errs := ValidateMap(body, validatorMap)
	if len(errs) > 0 {
		return nil, ServerError{400, Errors{Errors: errs}}
	}
	return body, nil
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
	} else {
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

