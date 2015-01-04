package rest

import (
	"github.com/KoFish/pallium/api"
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func setupProfile(root *mux.Router) {
	root.HandleFunc("/profile/{user}/displayname", u.OptionsReply()).Methods("OPTIONS")
	root.HandleFunc("/profile/{user}/avatar_url", u.OptionsReply()).Methods("OPTIONS")
	root.Handle("/profile/{user}/avatar_url", u.JSONReply(u.RequireAuth(getAvatarURL))).Methods("GET")
	root.Handle("/profile/{user}/avatar_url", u.JSONReply(u.RequireAuth(updateAvatarURL))).Methods("PUT")
	root.Handle("/profile/{user}/displayname", u.JSONReply(u.RequireAuth(getDisplayName))).Methods("GET")
	root.Handle("/profile/{user}/displayname", u.JSONReply(u.RequireAuth(updateDisplayName))).Methods("PUT")
}

func getDisplayName(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetDisplayName(user, r.Body, mux.Vars(r))
}

func updateDisplayName(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetDisplayName(user, r.Body, mux.Vars(r))
}

func getAvatarURL(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetDisplayName(user, r.Body, mux.Vars(r))
}

func updateAvatarURL(user *s.User, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	return api.GetDisplayName(user, r.Body, mux.Vars(r))
}
