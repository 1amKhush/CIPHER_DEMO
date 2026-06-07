package main

import (
	// "fmt"

	"net/http"

	"github.com/gorilla/mux"
)


// ServeLoginPage handles rendering the login GUI view
func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Make sure the file exists at this path relative to where you run your main binary
	http.ServeFile(w, r, "frontend/login.html")
}

// ServeRegisterPage handles rendering the user registration view
func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "frontend/register.html")
}

// ServeTasksPage handles rendering the core dashboard UI layout view
func ServeTasksPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "frontend/tasks.html")
}




func Route() *mux.Router {
	router := mux.NewRouter()

	// Unprotected Public Routes (Page serving)
	router.HandleFunc("/", ServeHome).Methods("GET") 
	router.HandleFunc("/register", ServeRegisterPage).Methods("GET")
	router.HandleFunc("/login", ServeLoginPage).Methods("GET")
	router.HandleFunc("/tasks", ServeTasksPage).Methods("GET")

	// Authentication API Routes
	router.HandleFunc("/register", RegisterUser).Methods("POST")
	router.HandleFunc("/login", LoginUser).Methods("POST")

	// Protected Profile Route
	router.Handle("/profile", AuthMiddleware(http.HandlerFunc(Profile))).Methods("GET")

	// Protected API task routes (prefixed with /api to avoid conflicts with HTML routes) // for postman testing 
	router.Handle("/api/tasks", AuthMiddleware(http.HandlerFunc(ServeTasks))).Methods("GET")
	router.Handle("/api/task/{id}", AuthMiddleware(http.HandlerFunc(ServeTask))).Methods("GET")
	router.Handle("/api/task", AuthMiddleware(http.HandlerFunc(addTask))).Methods("POST")
	router.Handle("/api/task/{id}", AuthMiddleware(http.HandlerFunc(updateTask))).Methods("PUT")
	router.Handle("/api/task/{id}", AuthMiddleware(http.HandlerFunc(deleteTask))).Methods("DELETE")

	return router
}
