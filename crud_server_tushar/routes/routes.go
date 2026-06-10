package routes

import (
	"crud_server/controllers"
	"crud_server/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Register_routes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.Homeroute).Methods("GET")             //done
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST") //done
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")       //done

	router.Handle("/notes", middleware.RequireAuth(http.HandlerFunc(controllers.GetAllNotes))).Methods("GET")       //done  // even though this doesnt need auth but still just kept it like that
	router.Handle("/note/{id}", middleware.RequireAuth(http.HandlerFunc(controllers.GetNote))).Methods("GET")       //done
	router.Handle("/note", middleware.RequireAuth(http.HandlerFunc(controllers.CreateNote))).Methods("POST")        //done
	router.Handle("/note/{id}", middleware.RequireAuth(http.HandlerFunc(controllers.UpdateNote))).Methods("PUT")    //done
	router.Handle("/note/{id}", middleware.RequireAuth(http.HandlerFunc(controllers.DeleteNote))).Methods("DELETE") //done

	return router
}
