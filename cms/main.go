package cms

import (
    "net/http"
    "github.com/go-martini/martini"
)

func init() {
    m := martini.Classic()

    m.Handlers(AuthenticationHandler) // binds appengine.Context

    m.Get("/item/?$", GetSlugsHandler)

    m.Get("/item/(?P<slug>[a-zA-Z\\-]+)", GetItemHandler)

    m.Post("/item/(?P<slug>[a-zA-Z\\-]+)/publish", PublishItemHandler)

    m.Post("/item/(?P<slug>[a-zA-Z\\-]+)", PostItemHandler)

    m.Delete("/item/(?P<slug>[a-zA-Z\\-]+)", DeleteItemHandler)

    http.Handle("/", m)
}