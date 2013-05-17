package main

import (
  "fmt"
  "io/ioutil"
)

type Page struct {
  Title string
  Body  []byte	// a byte "slice". we don't use 'string' as '[]byte' is expected by io libraries
}

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

func main() {
  p1 := &Page{Title: "TestPage", Body: []byte("This is a simple Page.")}
  p1.save()
  p2, _ := loadPage("TestPage")
  fmt.Println(string(p2.Body))
}