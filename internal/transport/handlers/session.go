package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	sessionKey []byte
	store      *sessions.CookieStore
)

func InitSessionStore(key string) {
	if key == "" {
		key = "change-me-in-env"
	}
	sessionKey = []byte(key)
	store = sessions.NewCookieStore(sessionKey)
}

func getUserIDFromSession(r *http.Request) (int, error) {
	if store == nil {
		return 0, http.ErrNoCookie
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		return 0, err
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0, http.ErrNoCookie
	}

	return userID, nil
}

func setUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	if store == nil {
		return http.ErrNoCookie
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = userID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
	}

	return session.Save(r, w)
}

func clearUserSession(w http.ResponseWriter, r *http.Request) error {
	if store == nil {
		return nil
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := clearUserSession(w, r); err != nil {
		http.Error(w, "failed to clear session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
