package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
)

type Page struct {
  Title string
  Body  []byte	// a byte "slice". we don't use 'string' as '[]byte' is expected by io libraries
}

const lenPath = len("/view/")

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

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[lenPath:]
  p, _ := loadPage(title)
  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[lenPath:]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  fmt.Fprintf(w, "<h1>Editing %s</h1>"+
                 "<form action=\"/save/%s\" method=\"POST\">"+
                 "<textarea name=\"body\">%s</textarea><br>"+
                 "<input type=\"submit\" value=\"Save\">"+
                 "</form>",
              p.Title, p.Title, p.Body)
}

//func main() {
//  p1 := &Page{Title: "TestPage", Body: []byte("This is a simple Page.")}
//  p1.save()
//  p2, _ := loadPage("TestPage")
//  fmt.Println(string(p2.Body))
//}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
//  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}