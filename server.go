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
  index := template.Must(template.New("layout.html").ParseFiles("templates/layout.html"))
  if err := index.Execute(w, map[string]interface{}{
    "User": u,
    "LoginOrOutUrl": loginOrOutUrl,
  }); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func ServeAdmin(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  index := template.Must(template.New("admin.html").ParseFiles("templates/admin.html"))
  if err := index.Execute(w, map[string]interface{}{}); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}