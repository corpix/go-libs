package xerrors

import (
	"fmt"
)

type PanicError struct {
	Payload any
}

func (e PanicError) Error() string {
	return fmt.Sprintf("panic() was raised with payload: %+v", e.Payload)
}

var _ error = &PanicError{}

// func CatchPanic(f func()) (err *PanicError) {
// 	defer func() {
// 		if p := recover(); p != nil {
// 			err = &PanicError{Payload: p}
// 		}
// 	}()

// 	f()

// 	return nil
// }

// func CatchPanicInErr(f func() error) (err error) {
// 	defer func() {
// 		if p := recover(); p != nil {
// 			err = PanicError{Payload: p}
// 		}
// 	}()

// 	return f()
// }

func Recover[T any](f func() (T, error)) (_ T, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = PanicError{Payload: p}
		}
	}()

	return f()
}

func RecoverVoid(f func() error) error {
	_, err := Recover(func() (any, error) {
		return nil, f()
	})

	return err
}

// func GoRecover(f func()) <-chan PanicError {
// 	out := make(chan PanicError, 1)

// 	go func() {
// 		defer close(out)
// 		p := CatchPanic(f)
// 		if p != nil {
// 			out <- *p
// 		}
// 	}()

// 	return out
// }

// func GoRecoverForever(f func()) <-chan PanicError {
// 	out := make(chan PanicError)

// 	go func() {
// 		defer close(out)
// 		for {
// 			p := CatchPanic(f)
// 			if p == nil {
// 				break
// 			}

// 			out <- *p
// 		}
// 	}()

// 	return out
// }
