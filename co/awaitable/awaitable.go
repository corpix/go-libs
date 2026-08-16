package awaitable

import (
	"context"
)

type Awaitable interface {
	Await(ctx context.Context) error
}
