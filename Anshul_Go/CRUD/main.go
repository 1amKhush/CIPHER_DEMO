package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Note struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User
var sessions = map[string]string{}

var notes []Note
var nextID = 1

func main() {

	loadUsers()
	loadNotes()

	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/", showNotes)
	http.HandleFunc("/create", createNote)
	http.HandleFunc("/edit", editNote)
	http.HandleFunc("/delete", deleteNote)
	fmt.Println("Running at http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

func loadUsers() {
	data, err := os.ReadFile("users.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, &users)
}

func saveUsers() {
	data, _ := json.MarshalIndent(users, "", "  ")

	os.WriteFile("users.json", data, 0644)
}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getCurrentUser(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	return sessions[cookie.Value]
}

func showNotes(w http.ResponseWriter, r *http.Request) {
	username := getCurrentUser(r)

	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/notes.html"))
	var userNotes []Note

	for _, note := range notes {
		if note.Username == username {
			userNotes = append(userNotes, note)
		}
	}

	tmpl.Execute(w, userNotes)
}

func register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if r.Method == "POST" {
		if username == "" || password == "" {
			http.Error(w, "Username and password required", http.StatusBadRequest)
			return
		}

		for _, user := range users {
			if user.Username == username {
				http.Error(w, "Username already exists", http.StatusBadRequest)
				return
			}
		}
		users = append(users, User{Username: username, Password: password})

		saveUsers()

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		for _, user := range users {
			if user.Username == username && user.Password == password {
				sessionID := generateSessionID()

				sessions[sessionID] = username
				http.SetCookie(w, &http.Cookie{Name: "session_id", Value: sessionID, Path: "/"})

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")

	if err == nil {
		delete(sessions, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	username := getCurrentUser(r)

	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		notes = append(notes, Note{ID: nextID, Username: username, Title: title, Content: content})

		nextID++

		saveNotes()
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
	username := getCurrentUser(r)

	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, "Invalid Note ID", http.StatusBadRequest)
		return
	}

	for i := range notes {
		if notes[i].ID == id && notes[i].Username == username {
			notes = append(notes[:i], notes[i+1:]...)

			saveNotes()
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Note not found or access denied", http.StatusForbidden)
}

func editNote(w http.ResponseWriter, r *http.Request) {
	username := getCurrentUser(r)
	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		http.Error(w, "Invalid Note ID", http.StatusBadRequest)
		return
	}

	for i := range notes {
		if notes[i].ID == id && notes[i].Username == username {
			if r.Method == "POST" {
				notes[i].Title = r.FormValue("title")
				notes[i].Content = r.FormValue("content")

				saveNotes()
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			tmpl := template.Must(template.ParseFiles("templates/edit.html"))
			tmpl.Execute(w, notes[i])
			return
		}
	}
	http.Error(w, "Note not found or Access denied", http.StatusForbidden)
}

func loadNotes() {
	data, err := os.ReadFile("notes.json")

	if err != nil {
		return
	}

	json.Unmarshal(data, &notes)

	for _, note := range notes {
		if note.ID >= nextID {
			nextID = note.ID + 1
		}
	}
}

func saveNotes() {
	data, err := json.MarshalIndent(notes, "", "  ")

	if err != nil {
		return
	}
	os.WriteFile("notes.json", data, 0644)
}
