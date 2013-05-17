package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
  "regexp"
  "errors"
)

type Page struct {
  Title string
  Body  []byte	// a byte "slice". we don't use 'string' as '[]byte' is expected by io libraries
}

const lenPath = len("/view/")

// The function template.Must is a convenience wrapper that panics when
// passed a non-nil error value, and otherwise returns the *Template unaltered
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

// This is a method named 'save' that takes as its receiver p, a pointer to Page.
// It take no parameters and returns a value of type 'error'.
func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
  // returns `nil` if everything ok
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {
  title = r.URL.Path[lenPath:]
  if !titleValidator.MatchString(title) {
    http.NotFound(w, r)
    err = errors.New("Invalid Page Title")
  }
  return
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title, err := getTitle(w, r)
  if err != nil {
    return
  }
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title, err := getTitle(w, r)
  if err != nil {
    return
  }
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title, err := getTitle(w, r)
  if err != nil {
    return
  }
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err = p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}