package cms

import (
    "appengine"
    "appengine/user"
    "net/http"
    "encoding/json"
    "errors"
    "github.com/go-martini/martini"
)


type LoginResponse struct {
    Url    string      `json:"url"`
}


var ErrUnauthorized = errors.New("auth: Unauthorized")


func Authenticate(w http.ResponseWriter, r *http.Request ) (appengine.Context, error) {

    c := appengine.NewContext(r)

    u := user.Current(c)
    if u == nil {
        url, _ := user.LoginURL(c, r.URL.String())
        response := LoginResponse{Url:url}

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        b, err := json.Marshal(response)
        if err != nil {
            return c, err
        }

        http.Error(w, string(b), 401)
        return c, ErrUnauthorized
    }

    return c, nil
}

func AuthenticationHandler(c martini.Context, w http.ResponseWriter, r *http.Request) {

    var ac appengine.Context
    err := ErrUnauthorized

    if r.Method == "GET" {
        ac = appengine.NewContext(r)
    } else {

        ac, err = Authenticate(w, r)

        if err == ErrUnauthorized {
            // error already written to response
            return
        }
    }

    c.Map(ac) // *appengine.Context
}