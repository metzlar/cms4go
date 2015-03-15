package cms

import (
    "appengine"
    "appengine/datastore"
    "time"
    "fmt"
    "encoding/json"
    "net/http"
    "io"
    "io/ioutil"
)


type Item struct {
    Content     string      `json:"content"`
    Slug        string      `json:"slug"`
    Archived    time.Time   `json:"-"`
    Published   time.Time   `json:"published"`

    // The google appengine datastore already allows for Hierarchical data
    // so no parent reference is necessary
    // ParentKey datastore.Key
}


func (i *Item) String() string {
    return fmt.Sprintf("<Item %s>", i.Slug)
}


func ItemKey(c appengine.Context, slug string) *datastore.Key{
    return datastore.NewKey(c, "Item", slug, 0, nil)
}


func GetItemBySlug(c appengine.Context, slug string) (*Item, error) {
    return _GetItemBySlug(c, slug, false, false)
}


func _GetItemBySlug(c appengine.Context, slug string, includeArchived bool, includeUnpublished bool) (*Item, error) {
    query := datastore.NewQuery("Item").Filter("Slug =", slug)
    if includeArchived == false {
        query = query.Filter("Archived <=", 1)
    }
    // Only one inequality filter per query is supported.
    // if includeUnpublished == false {
    //    query = query.Filter("Published <=", time.Now())
    // }
    query = query.Limit(1000)
    result := make([]Item, 0, 1000)
    if _, err := query.GetAll(c, &result); err != nil {
        return nil, err
    }

    result = FilterUnPublished(result)

    if len(result) == 0 {
        return nil, datastore.ErrNoSuchEntity
    }

    return &result[0], nil
}


func FilterUnPublished(set []Item) []Item {
    var tNow = time.Now()
    var filtered []Item // == nil
    for _, item := range set {
        if item.Published.Before(tNow) {
            filtered = append(filtered, item)
        }
    }
    return filtered
}


func GetItemSlugs(c appengine.Context) ([]string, error) {
    var tNow = time.Now()
    var filtered []string // == nil
    query := datastore.NewQuery("Item").Filter("Archived <=", 1).Limit(1000)
    for iter := query.Run(c); ; {
        var item Item
        _, err := iter.Next(&item)
        if err == datastore.Done {
            break
        }
        if err != nil {
            return nil, err
        }
        if(item.Published.Before(tNow)) {
            filtered = append(filtered, item.Slug)
        }
    }
    return filtered, nil
}


func SaveItem(c appengine.Context, item *Item) error {
    key := ItemKey(c, item.Slug)
    _, err := datastore.Put(c, key, item)
    if err != nil {
        return err
    }
    return nil
}

func SerializeItem(w http.ResponseWriter, item *Item) error {
    b, err := json.Marshal(item)

    if err != nil {
        return err
    }

    w.Write(b)

    return nil
}

func DeserializeItem(body io.ReadCloser, item *Item) error {
    b, err := ioutil.ReadAll(body)

    if err != nil {
        return err
    }

    return json.Unmarshal(b, &item)
}