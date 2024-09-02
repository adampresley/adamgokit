package session

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type Session[T any] interface {
	Get(r *http.Request) (T, error)
	GetSession(r *http.Request) (*sessions.Session, error)
	GetStore() sessions.Store
	Save(w http.ResponseWriter, r *http.Request) error
	Set(r *http.Request, value T) error
}

type SessionWrapper[T any] struct {
	KeyName     string
	SessionName string
	Store       sessions.Store
}

func NewSessionWrapper[T any](store sessions.Store, sessionName, keyName string) SessionWrapper[T] {
	return SessionWrapper[T]{
		KeyName:     keyName,
		SessionName: sessionName,
		Store:       store,
	}
}

func (sw SessionWrapper[T]) Get(r *http.Request) (T, error) {
	var (
		err     error
		session *sessions.Session
		empty   T
	)

	if session, err = sw.GetSession(r); err != nil {
		return empty, fmt.Errorf("could not get session in Get: %w", err)
	}

	result := session.Values[sw.KeyName]
	return any(result).(T), nil
}

func (sw SessionWrapper[T]) GetSession(r *http.Request) (*sessions.Session, error) {
	return sw.Store.Get(r, sw.SessionName)
}

func (sw SessionWrapper[T]) GetStore() sessions.Store {
	return sw.Store
}

func (sw SessionWrapper[T]) Save(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = sw.GetSession(r); err != nil {
		return fmt.Errorf("could not get session in Get: %w", err)
	}

	return session.Save(r, w)
}

func (sw SessionWrapper[T]) Set(r *http.Request, value T) error {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = sw.GetSession(r); err != nil {
		return fmt.Errorf("could not get session in Set: %w", err)
	}

	session.Values[sw.KeyName] = value
	return nil
}