package memory

import "context"

type Cleaner interface {
	Schedule(ctx context.Context)
	Stop()
}
