package main

import (
	"log"
	"net/http"
	"workmate/pkg/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/task", handlers.CreateTask).Methods("POST")
	router.HandleFunc("/task", handlers.DeleteTask).Methods("DELETE")
	router.HandleFunc("/task", handlers.GetTask).Methods("GET")
	router.HandleFunc("/task/all", handlers.GetAllTask).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
