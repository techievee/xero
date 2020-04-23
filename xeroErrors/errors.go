package xeroErrors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// Response header
const (
	Failed = "failed"
	Retry  = "retry"
)

const (
	defaultBadRequestMsg = "bad_request"
)

var (
	stdErr = log.New(os.Stderr, "", 0)
)

// Error is the high level struct that's going to be handled in recovery middleware
type Error struct {
	Status    string      `json:"status"`
	Err       string      `json:"error"`
	Message   interface{} `json:"message"`
	RequestID string      `json:"request_id"`
	Traceback []string    `json:"traceback"`

	Code  int       `json:"-"`
	Time  time.Time `json:"-"`
	Inner error     `json:"-"`
}

// Error makes it compatible with `error` interface.
func (e Error) Error() string {
	return fmt.Sprintf("code=%d, message=%v", e.Code, e.Message)
}

// String returns string.
func (e Error) String() string {
	return e.Error()
}

// Error object for custom error handling
func New(code int, desc string, status string, message ...interface{}) Error {
	desc = strings.Replace(desc, " ", "_", -1)
	desc = strings.ToLower(desc)

	e := Error{
		Code:   code,
		Status: status,
		Err:    fmt.Sprintf("%s.%s", "errors", desc),
		Time:   time.Now(),
	}
	e.Message = e.Err

	if len(message) != 0 {
		if message[0] != nil {
			switch _t := message[0].(type) {
			case string:
				e.Message = _t
			case error:
				e.Message = _t.Error()
			default:
				e.Message = fmt.Sprintf("%v", _t)
			}

			e.Traceback = TrimStacktrace(strings.Split(fmt.Sprintf("%+v\n", message[0]), "\n"))
		}
	} else {
		stack := make([]byte, 1024*8)
		stack = stack[:runtime.Stack(stack, false)]
		traceback := strings.Split(string(stack), "\n\t")
		for i := range traceback {
			v := strings.Split(traceback[i], " ")
			if len(v) != 0 {
				traceback[i] = v[0]
			}
		}

		e.Traceback = TrimStacktrace(traceback)
	}

	return e
}

// Function helps in limiting the stackstrace to 50 lines
// TrimStacktrace trim the path from stacktrace and limit it to 15 lines
func TrimStacktrace(traceback []string) []string {

	if len(traceback) > 50 {
		traceback = traceback[:50]
	}

	for i := range traceback {
		if n := strings.SplitN(traceback[i], "xero", 2); len(n) > 1 {
			traceback[i] = n[1]
		}
	}

	return traceback
}

// XeroUnexpectedGenericError
// returns the same error if it is Error type, otherwise return 500 error and ask for retry
func NewUnexpectedGenericError(message interface{}) Error {
	// check err is Error type
	if err, ok := message.(Error); ok {
		if err.Code == 0 {
			err.Code = http.StatusInternalServerError
		}
		return err
	}
	return New(http.StatusInternalServerError, "unexpected_error", Retry, message)
}

// xeroForbiddenError
// returns 403 access defined error
// Returned when the resource requested by the user is denied for the current context due to role permission
func XeroForbiddenError() Error {
	return New(http.StatusForbidden, "access_denied", Failed)
}

// XeroUnauthorizedError
// returns 401 authentication failed or not provided error
// User authentication failed wrong username and password
func XeroUnauthorizedError() Error {
	return New(http.StatusUnauthorized, "unauthorized", Failed)
}

// XeroNotFoundError
// returns 404 Not Found
// Requested resources were not found in the server
func XeroNotFoundError(resourceType string) Error {
	return New(http.StatusNotFound, resourceType+"_unavailable", Failed)
}

// XeroBadRequestError returns 400 bad request error.
// The interface contains array
// The first argument is description
// The optional second argument is the message.
func XeroBadRequestError(err ...interface{}) Error {
	desc := defaultBadRequestMsg

	if len(err) != 0 {
		if err[0] != nil && err[0] != "" {
			desc = fmt.Sprintf("%v", err[0])
		}
		// Remove first item
		err = append(err[:0], err[1:]...)
	}

	return New(http.StatusBadRequest, desc, Failed, err...)
}

// LogStdError logs an Error to stdError
func LogStdError(e Error) {

	l := map[string]interface{}{
		"code":       e.Code,
		"err":        e.Err,
		"message":    e.Message,
		"stacktrace": e.Traceback,
	}

	payload, _ := json.MarshalIndent(l, "", "  ")
	stdErr.Println(string(payload))
}

// CatchLog wraps an error with a stacktrace and logs it to stdError.
func CatchLog(e error, reqID string) {
	if e == nil {
		return
	}

	CatchErr(&e)

	logError := NewUnexpectedGenericError(e)

	LogStdError(logError)
}

// CatchErr wraps the referenced error with a stacktrace and throws a panic, if occurs
// suppressPanic will prevent a panic from being rethrown
func CatchErr(oerr *error, suppressPanic ...bool) {
	recover := recover()

	defer func() {
		if recover != nil && !(len(suppressPanic) > 0 && suppressPanic[0]) && oerr != nil {
			panic(*oerr)
		}
	}()

	// assign new value to oerr if recover != nil
	if recover != nil {
		if oerr == nil {
			oerr = new(error)
		}
		if err, ok := recover.(error); ok {
			*oerr = err
		}
	}
}
