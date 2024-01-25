package supervisor

import "context"

type Supervisor interface {
	Start(ctx context.Context)
}
