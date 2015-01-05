package rest

import (
	"github.com/KoFish/pallium/api"
	u "github.com/KoFish/pallium/rest/utils"
	// s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	// "net/http"
)

func setupProfile(root *mux.Router) {
	root.HandleFunc("/profile/{user}/displayname", u.OptionsReply).Methods("OPTIONS")
	root.HandleFunc("/profile/{user}/avatar_url", u.OptionsReply).Methods("OPTIONS")
	root.Handle("/profile/{user}/avatar_url", u.AuthAPIEndpoint(api.GetAvatarURL)).Methods("GET")
	root.Handle("/profile/{user}/avatar_url", u.AuthAPIEndpoint(api.UpdateAvatarURL)).Methods("PUT")
	root.Handle("/profile/{user}/displayname", u.AuthAPIEndpoint(api.GetDisplayName)).Methods("GET")
	root.Handle("/profile/{user}/displayname", u.AuthAPIEndpoint(api.UpdateDisplayName)).Methods("PUT")
}
