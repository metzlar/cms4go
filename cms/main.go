package cms

import (
    "net/http"
    "github.com/go-martini/martini"
    "appengine"
    "appengine/datastore"
    "time"
    "encoding/json"
)

func init() {
    m := martini.Classic()

    m.Get("/item/?$", func(w http.ResponseWriter, r *http.Request){

        c := appengine.NewContext(r)

        slugs, err := GetItemSlugs(c)

        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        b, err2 := json.Marshal(slugs)

        if err2 != nil {
            http.Error(w, err2.Error(), 500)
            return
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.Write(b)
    })

    m.Get("/item/(?P<slug>[a-zA-Z\\-]+)", func(w http.ResponseWriter, r *http.Request, params martini.Params){

        c := appengine.NewContext(r)

        item, err := GetItemBySlug(c, params["slug"])
        if err != nil {
            http.Error(w, err.Error(), 404)
            return
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        err2 := SerializeItem(w, item)

        if err2 != nil {
            http.Error(w, err2.Error(), 500)
            return
        }
    })

    m.Post("/item/(?P<slug>[a-zA-Z\\-]+)/publish", func(w http.ResponseWriter, r *http.Request, params martini.Params){

        c, err := Authenticate(w, r)

        if err == ErrUnauthorized {
            // error already written to response
            return
        } else if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        existing, err := GetItemBySlug(c, params["slug"])

        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        existing.Published = time.Now()

        err2 := SaveItem(c, existing)
        if err2 != nil {
            http.Error(w, err2.Error(), 500)
            return
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        err3 := SerializeItem(w, existing)

        if err3 != nil {
            http.Error(w, err3.Error(), 500)
            return
        }
    })

    m.Post("/item/(?P<slug>[a-zA-Z\\-]+)", func(w http.ResponseWriter, r *http.Request, params martini.Params){

        c, err := Authenticate(w, r)

        if err == ErrUnauthorized {
            // error already written to response
            return
        } else if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        existing, err := GetItemBySlug(c, params["slug"])

        if err == datastore.ErrNoSuchEntity{
            existing = &Item{Slug:params["slug"]}
        } else if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        err = DeserializeItem(r.Body, existing)

        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        err = SaveItem(c, existing)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        err = SerializeItem(w, existing)

        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
    })

    m.Delete("/item/(?P<slug>[a-zA-Z\\-]+)", func(w http.ResponseWriter, r *http.Request, params martini.Params){

        c, err := Authenticate(w, r)

        if err == ErrUnauthorized {
            // error already written to response
            return
        } else if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        existing, err := GetItemBySlug(c, params["slug"])

        if err != nil {
            http.Error(w, err.Error(), 404)
            return
        }

        existing.Archived = time.Now()

        err2 := SaveItem(c, existing)
        if err2 != nil {
            http.Error(w, err2.Error(), 500)
            return
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        err3 := SerializeItem(w, existing)

        if err3 != nil {
            http.Error(w, err3.Error(), 500)
            return
        }
    })

    http.Handle("/", m)
}