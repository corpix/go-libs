package xerrors

import (
	"context"
	"fmt"
)

type WrappedError struct {
	Base  error // Высокоуровневая ошибка (например, ErrUserNotFound)
	Cause error // Низкоуровневая причина (например, sql.ErrNoRows)
}

// Реализуем интерфейс error
func (e *WrappedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: %v", e.Base, e.Cause)
	}
	return e.Base.Error()
}

// Реализуем метод Unwrap, чтобы errors.Is проверял ОБЕ ошибки в цепочке
func (e *WrappedError) Unwrap() error {
	return e.Cause
}

// Метод Is позволяет errors.Is сопоставлять эту структуру с базовой ошибкой e.Base
func (e *WrappedError) Is(target error) bool {
	return e.Base == target
}

// Хелпер для удобного создания ошибки
func WithCause(base error, cause error) error {
	return &WrappedError{Base: base, Cause: cause}
}

func WithCauseFromContext(base error, ctx context.Context) error {
	return &WrappedError{Base: base, Cause: context.Cause(ctx)}
}

func FromContext(ctx context.Context) error {
	return &WrappedError{Base: ctx.Err(), Cause: context.Cause(ctx)}
}
