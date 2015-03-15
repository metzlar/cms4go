package cms

import (
    "appengine"
    "appengine/user"
    "net/http"
    "encoding/json"
    "errors"
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