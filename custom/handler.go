package custom

import (
	"fmt"
	"net/http"
	"reflect"
)

// Handler - type
type Handler interface {
	Invoke(w http.ResponseWriter, req *http.Request) error
}

// functionHandler is the generic function type
type functionHandler func(w http.ResponseWriter, req *http.Request) error

// Invoke calls the handler, and serializes the response.
// If the underlying handler returned an error, or an error occurs during serialization, error is returned.
func (handler functionHandler) Invoke(w http.ResponseWriter, req *http.Request) error {
	return handler(w, req)
}

func errorHandler(e error) functionHandler {
	return func(w http.ResponseWriter, req *http.Request) error {
		return e
	}
}

// NewHandler creates a base lambda handler from the given handler function. The
// returned Handler performs JSON serialization and deserialization, and
// delegates to the input handler function.  The handler function parameter must
// satisfy the rules documented by Start.  If handlerFunc is not a valid
// handler, the returned Handler simply reports the validation error.
func NewHandler(handlerFunc interface{}) Handler {
	if handlerFunc == nil {
		return errorHandler(fmt.Errorf("handler is nil"))
	}
	handler := reflect.ValueOf(handlerFunc)
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		return errorHandler(fmt.Errorf("handler kind %s is not %s", handlerType.Kind(), reflect.Func))
	}

	return functionHandler(func(w http.ResponseWriter, req *http.Request) error {
		response := handler.Call([]reflect.Value{reflect.ValueOf(req)})

		// convert return values into (interface{}, error)
		var err error
		if len(response) > 0 {
			if errVal, ok := response[len(response)-1].Interface().(error); ok {
				err = errVal
			}
		}

		return err
	})
}
