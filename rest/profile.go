package rest

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func setupProfile(root *mux.Router) {
	root.HandleFunc("/profile/{profile}/avatar_url", u.JSONWithAuthReply(getAvatarUrl)).Methods("GET")
	root.HandleFunc("/profile/{profile}/displayname", u.OptionsReply()).Methods("OPTIONS")
	root.HandleFunc("/profile/{profile}/displayname", u.JSONWithAuthReply(getDisplayName)).Methods("GET")
}

func getDisplayName(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var (
		resp struct {
			DisplayName string `json:"display_name"`
		}
	)

	profile, err := user.GetProfile(db)
	if err != nil {
		fmt.Println(err)
		resp.DisplayName = ""
	} else {
		resp.DisplayName = profile.DisplayName
	}

	return resp, nil
}

func getAvatarUrl(db *sql.DB, user *s.User, w http.ResponseWriter, r *http.Request) (interface{}, error) {

	var (
		resp struct {
			AvatarUrl string `json:"avatar_url"`
		}
	)
	profileId := mux.Vars(r)["profile"]
	resp.AvatarUrl = fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=mm", md5.Sum([]byte(profileId)))

	return resp, nil
}
