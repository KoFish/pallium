package rest

import (
	u "github.com/KoFish/pallium/rest/utils"
	s "github.com/KoFish/pallium/storage"
	"github.com/gorilla/mux"
	"net/http"
	"database/sql"
	"fmt"
	"crypto/md5"
)

func setupProfile(root *mux.Router) {
	root.HandleFunc("/profile/{profile}/avatar_url", u.JSONWithAuthReply(getAvatarUrl)).Methods("GET")
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
