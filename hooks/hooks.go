package hooks

import (
	"context"
	"fmt"

func Hooks(ctx context.Context, event *github.WebhookEvent) error {
	fmt.Println("Hooks")
	return nil
}

