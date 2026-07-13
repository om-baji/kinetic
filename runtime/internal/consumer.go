package internal

import (
	"context"

	"github.com/om-baji/kinetic/shared"
)

func Consumer() {
	ctx := context.Background()
	for {
		flow := shared.DequeWorkflow(ctx)
	}
}
