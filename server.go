package blog

import (
  "fmt"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"

  "appengine"
  "appengine/user"
)

var r *mux.Router = mux.NewRouter()

func init() {
  r.HandleFunc("/", ServeHome)
  r.HandleFunc("/admin", ServeAdmin)
  r.HandleFunc("/admin/post/new", ServeNewPost)
  r.HandleFunc("/admin/post/edit/{slug}", ServeEditPost).Name("EditPost")
  r.HandleFunc("/admin/post/{slug}", ServeUpdatePost).Methods("POST", "PUT")
  r.HandleFunc("/blog/{date:\\d{4}/\\d{2}/\\d{2}}/{slug}", ServePostHandler).
    Name("post")
  http.Handle("/", r)
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  u := user.Current(ctx)
  
  posts, err := getAllPosts(ctx)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  
  if err := runView(w, "templates/index.html", map[string]interface{}{
    "User": u,
    "Posts": posts,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServePostHandler(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  vars := mux.Vars(r)
  slug := vars["slug"]
  u := user.Current(ctx)

  post, err := getPostFromSlugString(ctx, slug)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if err := runView(w, "templates/show-post.html", map[string]interface{}{
    "Post": post,
    "User": u,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeAdmin(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  u := user.Current(ctx)

  posts, err := getAllPosts(ctx)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if err := runView(w, "templates/admin.html", map[string]interface{}{
    "User": u,
    "Posts": posts,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeNewPost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  post, err := newPost(ctx)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, fmt.Sprintf("/admin/post/edit/%s", post.Slug.StringID()), http.StatusFound)
}

func ServeEditPost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  vars := mux.Vars(r)
  slug := vars["slug"]

  post, err := getPostFromSlugString(ctx, slug)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return 
  }

  if err := runView(w, "templates/edit-post.html", map[string]interface{}{
    "Post": post,
    "PostId": post.Slug.StringID(),
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeUpdatePost(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  vars := mux.Vars(r)
  slug := vars["slug"]

  formPost := &Post{
    Title: r.FormValue("title"),
    Content: r.FormValue("content"),
  }
  if isDraft, isDraftErr := strconv.ParseBool(r.FormValue("draft")); isDraftErr != nil {
    formPost.IsDraft = isDraft
  }

  post, err := updatePost(ctx, slug, formPost)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, fmt.Sprintf("/admin/post/edit/%s", post.Slug.StringID()), http.StatusFound)
}