package blog

import (
  "fmt"
  "html/template"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "time"

  "github.com/gorilla/mux"

  "appengine"
  "appengine/datastore"
  "appengine/user"
)

var templateFunctions = template.FuncMap {
  "postToUrl": func (p Post) string {
    return p.PublishedDate.Format("/blog/2016/03/31")
  },
}

var r *mux.Router = mux.NewRouter()

func init() {
  r.HandleFunc("/", ServeHome)
  r.HandleFunc("/admin", ServeAdmin)
  r.HandleFunc("/admin/post/new", ServeNewPost)
  r.HandleFunc("/admin/post/edit/{slug}", ServeEditPost)
  r.HandleFunc("/admin/post/{slug}", ServeUpdatePost).Methods("POST", "PUT")
  r.HandleFunc("/blog/{date:\\d{4}/\\d{2}/\\d{2}}/{slug}", ServePostHandler).
    Name("post")
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

  query := datastore.NewQuery("Posts").Order("-PublishedDate").Limit(10)
  posts := make([]Post, 0, 10)
  if _, query_error := query.GetAll(ctx, &posts); query_error != nil {
    http.Error(w, query_error.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  index := template.Must(template.New("layout.html").Funcs(templateFunctions).ParseFiles(
    "templates/layout.html",
    "templates/index.html"))
  if err := index.ExecuteTemplate(w, "base", map[string]interface{}{
    "User": u,
    "LoginOrOutUrl": loginOrOutUrl,
    "Posts": posts,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServePostHandler(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  vars := mux.Vars(r)
  slug := vars["slug"]

  key := postKey(ctx, slug)
  post := &Post{}
  if err := datastore.Get(ctx, key, post); err != nil {

  }
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  index := template.Must(template.New("layout.html").ParseFiles(
    "templates/layout.html",
    "templates/show-post.html"))
  if err := index.ExecuteTemplate(w, "base", map[string]interface{}{
    "Post": post,
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
  Slug *datastore.Key
  Title string
  Content string
  PublishedDate time.Time
  EditDate time.Time
  IsDraft bool
}

func (p *Post) Url() *url.URL {
  url, err := r.Get("post").URL("date", p.PublishedDate.Format("2006/01/02"), "slug", p.Slug.StringID())
  if err != nil {
    panic(err)
  }
  return url
}

func hyphenizeTitle (title string) string {
  return strings.Replace(title, " ", "-", -1)
}

func postKey(c appengine.Context, slug string) *datastore.Key {
  return datastore.NewKey(c, "Posts", slug, 0, nil)
}

func createSlug(ctx appengine.Context, p *Post) *datastore.Key {
  if p.Slug != nil {
    return p.Slug
  }
  hyphenatedTitle := hyphenizeTitle(p.Title)
  return postKey(ctx, hyphenatedTitle)
}

func ServeNewPost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  post := Post{
    Title: fmt.Sprintf("Post %s", time.Now().Format("20060102030405")),
    Content: "Fill me in",
    PublishedDate: time.Now(),
    EditDate: time.Now(),
    IsDraft: true,
  }
  slug := createSlug(ctx, &post)
  post.Slug = slug

  _, err := datastore.Put(ctx, slug, &post)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, fmt.Sprintf("/admin/post/edit/%s", slug.StringID()), http.StatusFound)
}

func ServeEditPost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-type", "text/html; charset=utf-8")

  vars := mux.Vars(r)
  slug := vars["slug"]

  key := postKey(ctx, slug)
  post := &Post{Slug: key}
  if err := datastore.Get(ctx, key, post); err != nil {

  }

  index := template.Must(template.New("layout.html").ParseFiles(
    "templates/layout.html",
    "templates/edit-post.html"))
  if err := index.ExecuteTemplate(w, "base", map[string]interface{}{
    "Post": post,
    "PostId": key.StringID(),
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeUpdatePost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  vars := mux.Vars(r)
  slug := vars["slug"]

  key := postKey(ctx, slug)
  post := &Post{}
  if err := datastore.Get(ctx, key, post); err != nil {

  }

  if post.Title != r.FormValue("title") {
    datastore.Delete(ctx, key)
    post.Slug = nil
    post.Title = r.FormValue("title")
    post.Slug = createSlug(ctx, post)
  }
  post.Content = r.FormValue("content")
  if isDraft, isDraftErr := strconv.ParseBool(r.FormValue("draft")); isDraftErr == nil {
    post.IsDraft = isDraft
  }

  _, puterr := datastore.Put(ctx, post.Slug, post)
  if puterr != nil {
    http.Error(w, puterr.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, fmt.Sprintf("/admin/post/edit/%s", post.Slug.StringID()), http.StatusFound)
}