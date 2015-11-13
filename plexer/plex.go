package plexer

import (
	"errors"
)

var (
	MissingPublisher = errors.New("A publisher is needed")
)

type (
	BatchPublisher  func(MessageSet) error
	SinglePublisher func([]string, string) error
	MessageSet      struct {
		Tags     []string
		Messages []Message
	}
	Message struct {
		Tags []string
		Data string
	}

	Plexer struct {
		batch  map[string]BatchPublisher
		single map[string]SinglePublisher
	}
)

func NewPlexer() *Plexer {
	plex := &Plexer{
		batch:  make(map[string]BatchPublisher, 0),
		single: make(map[string]SinglePublisher, 0),
	}

	return plex
}

func (plex *Plexer) AddBatcher(name string, observer BatchPublisher) {
	plex.batch[name] = observer
}

func (plex *Plexer) RemoveBatcher(name string) {
	delete(plex.batch, name)
}

func (plex *Plexer) AddObserver(name string, observer SinglePublisher) {
	plex.single[name] = observer
}

func (plex *Plexer) RemoveObserver(name string) {
	delete(plex.single, name)
}

func (plex *Plexer) Publish(messages MessageSet) error {

	for _, observer := range plex.batch {
		go observer(messages)
	}

	for _, observer := range plex.single {
		for _, message := range messages.Messages {
			message.Tags = append(message.Tags, messages.Tags...)
			go observer(message.Tags, message.Data)
		}
	}
	return nil
}

func (plex *Plexer) PublishSingle(tags []string, data string) error {

	messages := MessageSet{
		Tags: []string{},
		Messages: []Message{
			Message{
				Tags: tags,
				Data: data,
			},
		},
	}
	return plex.Publish(messages)
}
