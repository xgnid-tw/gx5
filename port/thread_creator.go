package port

import "context"

type ThreadCreator interface {
	CreateThread(ctx context.Context, channelID string, name string, message string) (string, error)
}
