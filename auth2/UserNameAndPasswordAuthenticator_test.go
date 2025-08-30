package auth2_test

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adampresley/adamgokit/auth2"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	gob.Register(&TestUser{})
}

func TestUserNameAndPassword(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[*TestUser](sessionStore, sessionName, sessionKey)

	assert.NotNil(t, provider)
}

func TestUserNameAndPasswordWithOptions(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithContextKey("custom-context"),
		auth2.WithDebug(true),
		auth2.WithExcludedPaths([]string{"/health", "/metrics"}),
		auth2.WithRedirectURL("/login"),
	)

	assert.NotNil(t, provider)
}

func TestDestroySession(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](sessionStore, sessionName, sessionKey)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := provider.DestroySession(w, r)
	assert.NoError(t, err)
}

func TestSaveSession(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](sessionStore, sessionName, sessionKey)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	testUser := TestUser{
		ID:   1,
		Name: "Test User",
	}

	err := provider.SaveSession(w, r, testUser)
	assert.NoError(t, err)
}

func TestMiddlewareWithNoSession(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](sessionStore, sessionName, sessionKey)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	middleware.ServeHTTP(w, r)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddlewareWithRedirectURL(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithRedirectURL("/login"),
	)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	middleware.ServeHTTP(w, r)

	assert.False(t, called)
	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
}

func TestMiddlewareWithCustomResponder(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	customResponderCalled := false
	responderFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		customResponderCalled = true
		w.WriteHeader(http.StatusForbidden)
	}

	provider := auth2.UserNameAndPassword[TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithResponder(responderFunc),
	)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	middleware.ServeHTTP(w, r)

	assert.False(t, called)
	assert.True(t, customResponderCalled)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMiddlewareWithExcludedPathsExact(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithExcludedPaths([]string{"/health", "/metrics"}),
		auth2.WithExcludedPathsExact(true),
	)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	middleware.ServeHTTP(w, r)

	assert.True(t, called)
}

func TestMiddlewareWithExcludedPathsPrefix(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithExcludedPaths([]string{"/api/public"}),
		auth2.WithExcludedPathsExact(false),
	)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/public/users", nil)

	middleware.ServeHTTP(w, r)

	assert.True(t, called)
}

func TestMiddlewareWithValidSession(t *testing.T) {
	var (
		capturedError error
		receivedUser  *TestUser
	)

	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[*TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithContextKey("custom-user"),
		auth2.WithErrorFunc(func(w http.ResponseWriter, r *http.Request, err error) {
			if err != nil {
				capturedError = err
			}
		}),
	)

	testUser := &TestUser{
		ID:   1,
		Name: "Test User",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("In middleware\n")
		if user := r.Context().Value("custom-user"); user != nil {
			fmt.Printf("Received user: %+v\n", user)
			receivedUser = user.(*TestUser)
		}

		w.WriteHeader(http.StatusOK)
	})

	middleware := provider.Middleware(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := provider.SaveSession(w, r, testUser)
	assert.NoError(t, err)

	r = httptest.NewRequest(http.MethodGet, "/", nil)

	for _, cookie := range w.Result().Cookies() {
		r.AddCookie(cookie)
	}

	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, r)

	assert.NoError(t, capturedError)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testUser.ID, receivedUser.ID)
	assert.Equal(t, testUser.Name, receivedUser.Name)
}

func TestMiddlewareWithInvalidSessionValue(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	provider := auth2.UserNameAndPassword[TestUser](sessionStore, sessionName, sessionKey)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	session, _ := sessionStore.Get(r, sessionName)
	session.Values[sessionKey] = "invalid-type"
	session.Save(r, w)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	for _, cookie := range w.Result().Cookies() {
		r.AddCookie(cookie)
	}

	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, r)

	assert.False(t, called)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMiddlewareWithCustomErrorFunc(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))
	sessionName := "test-session"
	sessionKey := "user"

	customErrorCalled := false
	errorFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		customErrorCalled = true
		w.WriteHeader(http.StatusBadRequest)
	}

	provider := auth2.UserNameAndPassword[TestUser](
		sessionStore,
		sessionName,
		sessionKey,
		auth2.WithErrorFunc(errorFunc),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	session, _ := sessionStore.Get(r, sessionName)
	session.Values[sessionKey] = "invalid-type"
	session.Save(r, w)

	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	middleware := provider.Middleware(handler)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	for _, cookie := range w.Result().Cookies() {
		r.AddCookie(cookie)
	}

	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, r)

	assert.False(t, called)
	assert.True(t, customErrorCalled)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

