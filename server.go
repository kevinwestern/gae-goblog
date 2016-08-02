package blog

import (
  "fmt"
  "html/template"
  "net/http"
  "strconv"
  "time"

  "github.com/gorilla/mux"

  "appengine"
  "appengine/datastore"
  "appengine/user"
)

func init() {
  r := mux.NewRouter()
  r.HandleFunc("/", ServeHome)
  r.HandleFunc("/admin", ServeAdmin)
  r.HandleFunc("/admin/post/new", ServeNewPost)
  r.HandleFunc("/admin/post/edit/{id}", ServeEditPost)
  r.HandleFunc("/admin/post/{id}", ServeUpdatePost).Methods("POST", "PUT")
  http.Handle("/", r)
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  u := user.Current(ctx)
  loginOrOutUrl := ""
  if u == nil {
    loginOrOutUrl, _ = user.LoginURL(ctx, "/")
  } else {
    loginOrOutUrl, _ = user.LogoutURL(ctx, "/")
  }
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  index := template.Must(template.New("layout.html").ParseFiles(
    "templates/layout.html",
    "templates/index.html"))
  if err := index.ExecuteTemplate(w, "base", map[string]interface{}{
    "User": u,
    "LoginOrOutUrl": loginOrOutUrl,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeAdmin(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  u := user.Current(ctx)
  loginOrOutUrl := ""
  if u == nil {
    loginOrOutUrl, _ = user.LoginURL(ctx, "/admin")
  } else {
    loginOrOutUrl, _ = user.LogoutURL(ctx, "/admin")
  }
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  index := template.Must(template.New("layout.html").ParseFiles(
    "templates/layout.html",
    "templates/admin.html"))
  if err := index.ExecuteTemplate(w, "base", map[string]interface{}{
    "User": u,
    "LoginOrOutUrl": loginOrOutUrl,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

type Post struct {
  Title string
  Content string
  PublishedDate time.Time
  EditDate time.Time
  IsDraft bool
}

//func postKey(c appengine.Context) *datastore.Key {
//  return datastore.NewKey(c, "Posts", "default_posts", 0, nil)
//}

func ServeNewPost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  post := Post{
    Title: "New Post",
    Content: "Fill me in",
    PublishedDate: time.Now(),
    EditDate: time.Now(),
    IsDraft: true,
  }

  incomplete_key := datastore.NewIncompleteKey(ctx, "Posts", nil)
  key, err := datastore.Put(ctx, incomplete_key, &post)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, fmt.Sprintf("/admin/post/edit/%s", key.Encode()), http.StatusFound)
}

func ServeEditPost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-type", "text/html; charset=utf-8")

  vars := mux.Vars(r)
  id := vars["id"]

  key, err := datastore.DecodeKey(id)
  ctx.Infof("key is %s", key)
  if err != nil {
    // TODO(kevin): Do something
  }
  post := &Post{}
  ctx.Infof("About to fetch post")
  if err := datastore.Get(ctx, key, post); err != nil {

  }

  index := template.Must(template.New("layout.html").ParseFiles(
    "templates/layout.html",
    "templates/edit-post.html"))
  if err := index.ExecuteTemplate(w, "base", map[string]interface{}{
    "Post": post,
    "PostId": key.Encode(),
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeUpdatePost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)

  vars := mux.Vars(r)
  id := vars["id"]

  key, err := datastore.DecodeKey(id)
  ctx.Infof("key is %s", key)
  if err != nil {
    // TODO(kevin): Do something
  }
  post := &Post{}
  ctx.Infof("About to fetch post")
  if err := datastore.Get(ctx, key, post); err != nil {

  }

  post.Title = r.FormValue("title")
  post.Content = r.FormValue("content")
  post.IsDraft, err = strconv.ParseBool(r.FormValue("draft"))

  _, putterr := datastore.Put(ctx, key, post)
  if putterr != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, fmt.Sprintf("/admin/post/edit/%s", key.Encode()), http.StatusFound)
}