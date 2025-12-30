package ports

import "github.com/edlingao/internal/blog/core"

type CommentsChannelRepo interface {
	Start()
	AddSubscriber(client *core.Client)
	UnSubscribe(client *core.Client)
	BroadcastEvent(event *core.Event)
}
