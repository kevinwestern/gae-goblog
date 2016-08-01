package blog

import (
  "html/template"
  "net/http"

  "github.com/gorilla/mux"

  "appengine"
  "appengine/user"
)

func init() {
  r := mux.NewRouter()
  r.HandleFunc("/", ServeHome)
  r.HandleFunc("/admin", ServeAdmin)
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