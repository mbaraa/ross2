package context

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/mbaraa/ross2/models"
)

// HandlerContext holds SessionHandlerFunc required attributes for a handler function
type HandlerContext struct {
	Res  http.ResponseWriter
	Req  *http.Request
	Sess models.Session
}

// ReadJSON reads the request body as json into the given target and returns an error with what went wrong
// it might fail if the target isn't a pointer, or if the body isn't a proper JSON
func (c *HandlerContext) ReadJSON(target interface{}) error {
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		c.Res.WriteHeader(http.StatusBadRequest)
		return errors.New("target not pointer")
	}

	err := json.NewDecoder(c.Req.Body).Decode(target)
	_ = c.Req.Body.Close()

	if err != nil {
		c.Res.WriteHeader(http.StatusBadRequest)
		return err
	}

	return nil
}

// WriteJSON writes the given object into the http response
func (c *HandlerContext) WriteJSON(src interface{}, failStatus int) error {
	err := json.NewEncoder(c.Res).Encode(src)
	if err != nil {
		c.Res.WriteHeader(failStatus)
		return err
	}
	return nil
}
