package sessionmanager

import (
	"errors"

	"github.com/gorilla/securecookie"
	"gopkg.in/boj/redistore.v1"
)

const (
	port               = ":6379"
	protocol           = "tcp"
	maxIdleConnections = 5
)

//StreamingSession testing abstraction
type StreamingSession struct {
	UserID    uint
	Username  string
	StreamURL string
	ImagePath string
}

//SessionManager testing abstraction
type SessionManager struct {
	store *redistore.RediStore
}

//NewSessionStore ok
func NewSessionStore() (*redistore.RediStore, error) {
	key := securecookie.GenerateRandomKey(32)
	return redistore.NewRediStore(maxIdleConnections, protocol, port, "", key)
}

//NewSessionManager ok
func NewSessionManager() (*SessionManager, error) {
	manager := &SessionManager{
		store: nil,
	}

	if manager == nil {
		err := errors.New("Could not allocate session manager")
		return nil, err
	}

	store, err := NewSessionStore()
	if err != nil {
		return nil, err
	}

	manager.store = store

	return manager, nil
}
