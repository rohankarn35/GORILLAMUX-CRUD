package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	routes := mux.NewRouter()
	s := routes.PathPrefix("/api").Subrouter()
	s.HandleFunc("/createProfile", createProfile).Methods("POST")
	// s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	s.HandleFunc("/getUserProfile", getUserProfile).Methods("POST")
	s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", s))
}
