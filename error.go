package httputils



import (
	"fmt"
	"net/http"
)

type Errors struct {
	Errors []Error `json:"errors"`
}

type ServerError struct {
	StatusCode int
	Errors     Errors
}

func (self ServerError) Error() string {
	return self.Errors.Error()
}

func (self ServerError) write(w http.ResponseWriter) {
	JSON(w, self.Errors, self.StatusCode)
}

func raise500(w http.ResponseWriter) {
	UndefinedKeyError("INTERNAL_SERVER_ERROR", "Something went wrong").writeWithCode(500, w)
}

func HTTP400() ServerError {
	return ServerError{400, Errors{[]Error{UndefinedKeyError("INVALID_REQUEST", "Invalid request")}}}
}

func HTTP401() ServerError {
	return ServerError{401, Errors{[]Error{UndefinedKeyError("UNAUTHORIZED", "Unauthorized user")}}}
}

func HTTP403() ServerError {
	return ServerError{403, Errors{[]Error{UndefinedKeyError("PERMISSION_DENIED", "Permission denied")}}}
}

func HTTP404() ServerError {
	return ServerError{404, Errors{[]Error{UndefinedKeyError("NOT_FOUND", "Item not found")}}}
}

type Error struct {
	Key         string   `json:"key"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
	Args        []string `json:"args, omitempty"`
}

func (self Error) writeWithCode(code int, w http.ResponseWriter) {
	ServerError{code, Errors{[]Error{self}}}.write(w)
}

func UndefinedKeyError(code string, description string) Error {
	return Error{"undefined", description, code, nil}
}

func (self Error) Error() string {
	return self.Code
}

func (self Errors) Error() string {
	return fmt.Sprintf("Occured %d errors", len(self.Errors))
}
