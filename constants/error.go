package constants

import (
	"fmt"
	"strings"
)

type LogicError struct {
	err     error
	Message string
	Code    int64
}

func (err LogicError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.Message, err.err)
	}
	return err.Message
}

func (err LogicError) wrap(inner error) error {
	return LogicError{Message: err.Message, err: inner}
}

func (err LogicError) Unwrap() error {
	return err.err
}

func (err LogicError) Is(target error) bool {
	ts := target.Error()
	return ts == err.Message || strings.HasPrefix(ts, err.Message+": ")
}

var (
	ErrNotImpl      = LogicError{Message: "not implemented", Code: 500}
	ErrOk           = LogicError{Message: "ok", Code: 200}
	ErrConfig       = LogicError{Message: "invalid config", Code: 400}
	ErrUnauthorized = LogicError{Message: "no authorization headers", Code: 401}
)
