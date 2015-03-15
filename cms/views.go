package cms

import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "encoding/json"
    "github.com/go-martini/martini"
    "time"
)


func GetSlugsHandler(
    c appengine.Context,
    w http.ResponseWriter,
    r *http.Request,
) {

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
}

func GetItemHandler(
    c appengine.Context,
    w http.ResponseWriter,
    r *http.Request,
    params martini.Params,
) {
    existing, err := GetItemBySlug(c, params["slug"])

    if err == datastore.ErrNoSuchEntity {
        http.Error(w, err.Error(), 404)
        return
    }

    err = SerializeItem(w, existing)

    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}


func PostItemHandler(
    c appengine.Context,
    w http.ResponseWriter,
    r *http.Request,
    params martini.Params,
){
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
}


func PublishItemHandler(
    c appengine.Context,
    w http.ResponseWriter,
    r *http.Request,
    params martini.Params,
){
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
}


func DeleteItemHandler(
    c appengine.Context,
    w http.ResponseWriter,
    r *http.Request,
    params martini.Params,
) {
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
}