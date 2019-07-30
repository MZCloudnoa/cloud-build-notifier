package notifier

import "context"

// HandlePubSub HandlePubSub
func HandlePubSub(ctx context.Context, m PubSubMessage) error {
	return NewNotifier(nil).HandlePubSub(m.Data)
}
