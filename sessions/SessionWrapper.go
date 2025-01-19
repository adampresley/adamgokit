package sessions

import (
	"fmt"
	"net/http"

	gorillasessions "github.com/gorilla/sessions"
)

var (
	ErrSessionKeyNotFound = fmt.Errorf("session key not found")
)

type Session[T any] interface {
	Get(r *http.Request) (T, error)
	GetSession(r *http.Request) (*gorillasessions.Session, error)
	GetStore() gorillasessions.Store
	Save(w http.ResponseWriter, r *http.Request) error
	Set(r *http.Request, value T) error
	Destroy(w http.ResponseWriter, r *http.Request) error
}

type SessionWrapper[T any] struct {
	KeyName     string
	SessionName string
	Store       gorillasessions.Store
}

func NewSessionWrapper[T any](store gorillasessions.Store, sessionName, keyName string) SessionWrapper[T] {
	return SessionWrapper[T]{
		KeyName:     keyName,
		SessionName: sessionName,
		Store:       store,
	}
}

func (sw SessionWrapper[T]) Get(r *http.Request) (T, error) {
	var (
		err     error
		session *gorillasessions.Session
		empty   T
	)

	if session, err = sw.GetSession(r); err != nil {
		return empty, fmt.Errorf("could not get session in Get: %w", err)
	}

	if _, ok := session.Values[sw.KeyName]; !ok {
		return empty, ErrSessionKeyNotFound
	}

	result := session.Values[sw.KeyName]
	return any(result).(T), nil
}

func (sw SessionWrapper[T]) GetSession(r *http.Request) (*gorillasessions.Session, error) {
	return sw.Store.Get(r, sw.SessionName)
}

func (sw SessionWrapper[T]) GetStore() gorillasessions.Store {
	return sw.Store
}

func (sw SessionWrapper[T]) Save(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		session *gorillasessions.Session
	)

	if session, err = sw.GetSession(r); err != nil {
		return fmt.Errorf("could not get session in Get: %w", err)
	}

	return session.Save(r, w)
}

func (sw SessionWrapper[T]) Set(r *http.Request, value T) error {
	var (
		err     error
		session *gorillasessions.Session
	)

	if session, err = sw.GetSession(r); err != nil {
		return fmt.Errorf("could not get session in Set: %w", err)
	}

	session.Values[sw.KeyName] = value
	return nil
}

func (sw SessionWrapper[T]) Destroy(w http.ResponseWriter, r *http.Request) error {
	var (
		err     error
		session *gorillasessions.Session
	)

	if session, err = sw.GetSession(r); err != nil {
		return fmt.Errorf("could not get session in Destroy: %w", err)
	}

	session.Options.MaxAge = -1
	session.Values[sw.KeyName] = new(T)

	return session.Save(r, w)
}
