package rest

import (
    "database/sql"
    "encoding/json"
    "fmt"
    c "github.com/KoFish/pallium/config"
    m "github.com/KoFish/pallium/matrix"
    u "github.com/KoFish/pallium/rest/utils"
    s "github.com/KoFish/pallium/storage"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

type LoginFlow struct {
    Type   string   `json:"type"`
    Stages []string `json:"stages,omitempty"`
}

type LoginFlowInformation struct {
    Flows []LoginFlow `json:"flows"`
}

type PasswordLoginRequest struct {
    Type     string `json:"type"`
    User     string `json:"user"`
    Password string `json:"password"`
}

type LoginResponse struct {
    UserID string `json:"user_id"`
    Token  string `json:"access_token"`
}

func getLoginInfo(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    flows := LoginFlowInformation{
        Flows: []LoginFlow{
            LoginFlow{
                Type:   "m.login.password",
                Stages: []string{},
            },
        },
    }
    return flows, nil
}

func (regInfo *PasswordLoginRequest) submitRegistrationRequest(db *sql.DB, w http.ResponseWriter, r *http.Request) (*LoginResponse, error) {
    username := regInfo.User
    password := regInfo.Password

    txn, err := db.Begin()
    if err != nil {
        fmt.Errorf("Could not create new db transaction: %v", err)
        return nil, u.NewError(m.M_FORBIDDEN, "Could not create a new transaction")
    }

    user_id, _ := m.ParseUserID(username) // Get a struct to compare hostpart
    if err != nil {
        return nil, u.NewError(m.M_FORBIDDEN, "Invalid User ID: "+err.Error())
    }

    fmt.Printf("Request to register user %s\n", user_id)
    if !user_id.IsMine() {
        return nil, u.NewError(m.M_FORBIDDEN, "Can not register user on other homeserver "+user_id.Domain())
    }
    err = s.CreateUser(txn, user_id)
    if err != nil {
        txn.Rollback()
        return nil, u.NewError(m.M_USER_IN_USE, "Username already taken: "+err.Error())
    }
    user, err := s.GetUser(txn, user_id)
    if err != nil {
        txn.Rollback()
        return nil, u.NewError(m.M_NOT_FOUND, "Could not fetch user: "+err.Error())
    }
    if err = user.SetPassword(txn, password); err != nil {
        txn.Rollback()
        return nil, u.NewError(m.M_NOT_FOUND, "Could not set password: "+err.Error())
    }
    token, err := user.GetAccessToken(txn)
    if err != nil {
        txn.Rollback()
        return nil, u.NewError(m.M_NOT_FOUND, "Could not get access token for user: "+err.Error())
    }
    err = txn.Commit()
    return &LoginResponse{user_id.String(), string(token.Token)}, err
}

func (loginInfo *PasswordLoginRequest) submitLoginPassword(db *sql.DB, w http.ResponseWriter, r *http.Request) (*LoginResponse, error) {
    username := loginInfo.User
    password := loginInfo.Password

    user_id, err := m.ParseUserID(username)

    if err != nil {
        return nil, u.NewError(m.M_FORBIDDEN, "Invalid User ID: "+err.Error())
    }

    if user_id.Domain() != c.Hostname {
        return nil, u.NewError(m.M_FORBIDDEN, "Can not register user namespaced to other host.")
    }

    fmt.Printf("Request to login user %s\n", user_id)
    user, err := s.GetUser(db, user_id)
    if err != nil {
        return nil, u.NewError(m.M_NOT_FOUND, "Could not get such a user: "+err.Error())
    }
    if user.Password.Equal(password) {
        token, err := user.GetAccessToken(db)
        if err != nil {
            return nil, u.NewError(m.M_NOT_FOUND, "Could not find a valid access token: "+err.Error())
        }
        if err = token.UpdateAccessToken(db); err != nil {
            return nil, u.NewError(m.M_FORBIDDEN, "Could not update access time of token: "+err.Error())
        }
        return &LoginResponse{user_id.String(), string(token.Token)}, err
    } else {
        return nil, u.NewError(m.M_FORBIDDEN, "Invalid password")
    }
}

func submitLogin(db *sql.DB, w http.ResponseWriter, r *http.Request) (interface{}, error) {
    if !u.CheckTxnId(r) {
        return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
    }
    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        return nil, err
    }

    var msg struct {
        Type string `json:"type"`
    }
    if err := json.Unmarshal(body, &msg); err != nil {
        return nil, err
    }

    switch msg.Type {
    case "m.login.password":
        var loginreq PasswordLoginRequest
        if err := json.Unmarshal(body, &loginreq); err != nil {
            return nil, err
        }
        if response, err := loginreq.submitLoginPassword(db, w, r); err != nil {
            return nil, err
        } else {
            return response, nil
        }
    default:
        return nil, u.NewError(m.M_NOT_FOUND, "Invalid login request type")
    }
}

func submitRegister(db *sql.DB, w http.ResponseWriter, r *http.Request) (interface{}, error) {
    if !u.CheckTxnId(r) {
        return nil, u.NewError(m.M_FORBIDDEN, "Request has already been sent")
    }
    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        return nil, err
    }

    var msg struct {
        Type string `json:"type"`
    }
    if err := json.Unmarshal(body, &msg); err != nil {
        return nil, u.NewError(m.M_NOT_JSON, "Could not parse json: "+err.Error())
    }

    switch msg.Type {
    case "m.login.password":
        var regreq PasswordLoginRequest
        if err := json.Unmarshal(body, &regreq); err != nil {
            return nil, u.NewError(m.M_BAD_JSON, "Did not contain requested information: "+err.Error())
        }
        if response, err := regreq.submitRegistrationRequest(db, w, r); err != nil {
            return nil, err
        } else {
            return response, nil
        }
    default:
        return nil, u.NewError(m.M_NOT_FOUND, "Invalid registration request type")
    }

}

func fallbackLogin(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Not implemented")
}

func setupLogin(root *mux.Router) {
    root.HandleFunc("/login", u.AllowMatrixOrg).Methods("OPTIONS")
    root.HandleFunc("/login", u.JSONReplyHandler(getLoginInfo)).Methods("GET")
    root.HandleFunc("/login", u.JSONWithDBReply(submitLogin)).Methods("POST")
    root.HandleFunc("/login/{txnId:[0-9]+}", u.JSONWithDBReply(submitLogin)).Methods("PUT").Name("loginPUT")
    root.HandleFunc("/login/fallback", fallbackLogin).Methods("GET")
}

func setupRegister(root *mux.Router) {
    // TODO: Change the following function if the registration and login paths
    //       differ.
    root.HandleFunc("/register", u.AllowMatrixOrg).Methods("OPTIONS")
    root.HandleFunc("/register", u.JSONReplyHandler(getLoginInfo)).Methods("GET")
    root.HandleFunc("/register", u.JSONWithDBReply(submitRegister)).Methods("POST")
    root.HandleFunc("/register/{txnId:[0-9]+}", u.JSONWithDBReply(submitRegister)).Methods("PUT")
}
