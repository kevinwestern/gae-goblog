package blog

import (
  "html/template"
  "net/http"
)

var templates map[string]*template.Template

func init() {
  if templates == nil {
    templates = make(map[string]*template.Template)
  }
  templates["templates/admin.html"] = template.Must(template.ParseFiles("templates/admin.html", "templates/layout.html"))
  templates["templates/edit-post.html"] = template.Must(template.ParseFiles("templates/edit-post.html", "templates/layout.html"))
  templates["templates/index.html"] = template.Must(template.ParseFiles("templates/index.html", "templates/layout.html"))
  templates["templates/show-post.html"] = template.Must(template.ParseFiles("templates/show-post.html", "templates/layout.html"))
}

func runView(w http.ResponseWriter, viewName string, data interface{}) error {
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  return templates[viewName].ExecuteTemplate(w, "base", data)
}