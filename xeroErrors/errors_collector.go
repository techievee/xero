package xeroErrors

import (
	"fmt"
	"reflect"
	"sync"
)

// DuplicatationOptions  can be configured to allow duplicates.
type DuplicatationOptions int

const (
	AllowDuplicates  DuplicatationOptions = 0
	RejectDuplicates DuplicatationOptions = 1
)

// DefaultErrorFormatter represents the default formatter for displaying the collection
// of errors.
var DefaultErrorFormatter = func(i int, n int, err error, str *string) {
	*str = fmt.Sprintf("%s\n%d:%s", *str, i, err.Error())
}

// ErrorCollection allows multiple errors to be accumulated and then returned as a single error.
// ErrorCollection can be safely used by concurrent goroutines.
type ErrorCollection struct {
	DuplicatationOptions DuplicatationOptions
	Errors               []error
	Formatter            func(i int, n int, err error, str *string)
	lock                 sync.RWMutex
}

// NewErrorCollection creates a new empty ErrorCollection.
// When `dup` is set, any duplicate error message is discarded
// and not appended to the collection
func NewErrorCollection(dup ...DuplicatationOptions) *ErrorCollection {
	ec := &ErrorCollection{}
	ec.Errors = []error{}
	ec.Formatter = DefaultErrorFormatter
	if len(dup) != 0 {
		ec.DuplicatationOptions = dup[0]
	}
	return ec
}

// Append an error to the error collection without locking
func (ec *ErrorCollection) addError(err error) {

	if err == nil {
		return
	}

	if ec.DuplicatationOptions != AllowDuplicates {
		//Don't append if err is a duplicate
		for _, containedErr := range ec.Errors {
			//Use Reflection
			if reflect.DeepEqual(containedErr, err) {
				return
			}
		}
	}
	ec.Errors = append(ec.Errors, err)
}

// AddError appends an error to the error collection.
// It is safe to use from multiple concurrent goroutines.
func (ec *ErrorCollection) AddError(err error) {
	ec.lock.Lock()
	defer ec.lock.Unlock()

	ec.addError(err)
}

// AddErrors appends multiple errors to the error collection.
// It is safe to use from multiple concurrent goroutines.
func (ec *ErrorCollection) AddErrors(errs ...error) {
	ec.lock.Lock()
	defer ec.lock.Unlock()

	for _, err := range errs {
		ec.addError(err)
	}
}

// AddErrorCollection appends an entire ErrorCollection to the receiver error collection.
// It is safe to use from multiple concurrent goroutines.
func (ec *ErrorCollection) AddErrorCollection(errs *ErrorCollection) {
	ec.lock.Lock()
	defer ec.lock.Unlock()

	for _, err := range errs.Errors {
		ec.addError(err)
	}
}

// Error return a list of all contained errors.
// The output can be formatted by setting a custom Formatter.
func (ec *ErrorCollection) Error() string {
	if ec.Formatter == nil {
		return ""
	}

	ec.lock.RLock()
	defer ec.lock.RUnlock()
	str := ""
	for i, err := range ec.Errors {
		if ec.Formatter != nil {
			ec.Formatter(i, len(ec.Errors), err, &str)
		}
	}
	return str
}

// IsNil returns whether an error is nil or not.
// It can be used with ErrorCollection or generic errors
func IsNil(err error) bool {
	switch v := err.(type) {
	case *ErrorCollection:
		if len(v.Errors) == 0 {
			return true
		}
		return false

	default:
		if err == nil {
			return true
		}
		return false
	}
}
