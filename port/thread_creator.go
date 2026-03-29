package port

import "context"

type ThreadCreator interface {
	CreateThread(ctx context.Context, channelID string, name string) (string, error)
	SendThreadMessage(ctx context.Context, threadID string, message string) error
}
