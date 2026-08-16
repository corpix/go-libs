package xerrors

import (
	"fmt"
)

type FieldError struct {
	msg    string
	fields map[string]any
}

func (err *FieldError) Error() string {
	return err.msg
}

func (err *FieldError) Fields() map[string]any {
	return err.fields
}

func NewF(msg string, args ...any) error {
	pairCount := len(args) / 2

	err := FieldError{
		msg:    msg,
		fields: make(map[string]any, pairCount),
	}

	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]

		stringKey := fmt.Sprint(key)

		err.fields[stringKey] = value
	}

	return &err
}
