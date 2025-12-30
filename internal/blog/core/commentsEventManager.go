package core

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/labstack/echo/v4"
)

type Client struct {
	ID      string
	Title   string
	writter *echo.Response
}

type Event struct {
	Name  string
	Title string
	Data  map[string]any
}

type CommentsEventManager struct {
	Subscribers []*Client
	Subscribe   chan *Client
	Unsubscribe chan *Client
	Events      chan *Event
}

func NewCommentsEventManager() *CommentsEventManager {
	return &CommentsEventManager{
		Subscribers: []*Client{},
		Subscribe:   make(chan *Client),
		Unsubscribe: make(chan *Client),
		Events:      make(chan *Event),
	}
}

func (manager *CommentsEventManager) Start() {
	log.Println("Comments Event Manager started")
	for {
		select {
		case client := <-manager.Subscribe:
			manager.AddSubscriber(client)
		case client := <-manager.Unsubscribe:
			manager.UnSubscribe(client)
		case event := <-manager.Events:
			manager.BroadcastEvent(event)
		}
	}
}

func (manager *CommentsEventManager) AddSubscriber(client *Client) {
	manager.Subscribers = append(manager.Subscribers, client)
}

func (manager *CommentsEventManager) UnSubscribe(client *Client) {
	manager.Subscribers = slices.DeleteFunc(manager.Subscribers, func(c *Client) bool {
		return c.ID == client.ID
	})
}

func (manager *CommentsEventManager) BroadcastEvent(event *Event) {
	for _, client := range manager.Subscribers {
		if client.Title != event.Title {
			continue
		}
		err := client.SendEvent(event)
		if err != nil {
			manager.Subscribers = slices.DeleteFunc(manager.Subscribers, func(c *Client) bool {
				return c.ID == client.ID
			})
		}
	}
}

func NewClient(id, title string, writter *echo.Response) *Client {
	return &Client{
		ID:      id,
		Title:   title,
		writter: writter,
	}
}

func (client *Client) SendEvent(event *Event) error {
	_, err := client.writter.Write([]byte("event: " + event.Name + "\n"))
	if err != nil {
		return err
	}

	data, err := event.JSON()
	if err != nil {
		log.Println("Error marshalling event data:", err)
		return err
	}

	_, err = client.writter.Write([]byte("data: " + data + "\n\n"))
	if err != nil {
		log.Println("Error writing to client:", err)
		return err
	}

	client.writter.Flush()

	return nil
}

func NewEvent(name, title string, data map[string]any) *Event {
	return &Event{
		Name:  name,
		Title: title,
		Data:  data,
	}
}

func (event *Event) JSON() (string, error) {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
