package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Note struct {
	ID      int
	Title   string
	Content string
}

var notes []Note
var nextID = 1

func main() {
	http.HandleFunc("/", basicAuth(DisplayNotes))
	http.HandleFunc("/create", basicAuth(createNote))
	http.HandleFunc("/edit", basicAuth(editNote))
	http.HandleFunc("/delete", basicAuth(deleteNote))

	fmt.Println("Server running on http://localhost:9090")

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err)
	}
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		username, password, ok := r.BasicAuth()

		fmt.Println("username =", username)
		fmt.Println("password =", password)
		fmt.Println("ok =", ok)

		if !ok || username != "admin" || password != "1234" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func DisplayNotes(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/notes.html"))
	tmpl.Execute(w, notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		title := r.FormValue("title")
		content := r.FormValue("content")

		notes = append(notes, Note{
			ID:      nextID,
			Title:   title,
			Content: content,
		})

		nextID++

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/create.html")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, nil)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	for i := range notes {

		if notes[i].ID == id {

			notes = append(notes[:i], notes[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func editNote(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	for i := range notes {

		if notes[i].ID == id {

			if r.Method == "POST" {

				notes[i].Title = r.FormValue("title")
				notes[i].Content = r.FormValue("content")

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			tmpl := template.Must(
				template.ParseFiles("templates/edit.html"),
			)

			tmpl.Execute(w, notes[i])
			return
		}
	}
}
