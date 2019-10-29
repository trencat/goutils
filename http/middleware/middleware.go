package middleware

import (
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// AssertMethod returns 405 Method Not Allowed to requests not matching
// the given method.
func AssertMethod(method string) Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.ToUpper(r.Method) != method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("405 Method Not Allowed"))
			}
			fn(w, r)
		}
	}
}

func logRequest(r *http.Request, log *syslog.Writer) error {
	msg := fmt.Sprintf("%s %s %s%s %s %v %d",
		r.Proto, strings.ToUpper(r.Method), r.Host, r.RequestURI, r.RemoteAddr, r.Header, r.ContentLength)

	if r.ContentLength > 0 {
		// Read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err := errors.Errorf("%s Cannot read body", msg)
			log.Warning(fmt.Sprintf("%+v", err))
			return err
		}

		msg += fmt.Sprintf("%s %s", msg, body)
	}

	log.Debug(msg) // Log request
	return nil
}

// Log returns a middleware that logs requests to syslog
func Log(log *syslog.Writer) Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logRequest(r, log)
			fn(w, r)
			// Log response?
		}
	}
}
