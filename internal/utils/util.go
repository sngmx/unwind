package utils

import (
	"net/http"
	"unwind/internal/types"

	"github.com/gorilla/sessions"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request, store sessions.Store) (*types.UserInfo, error) {
	session, err := store.Get(r, "unwind-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Username not found in session", http.StatusInternalServerError)
	}

	return &types.UserInfo{
		Username: username,
	}, nil
}
